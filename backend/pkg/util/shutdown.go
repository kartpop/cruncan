package util

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type TerminatorFunc func(ctx context.Context) error

func GracefulShutdown(server *http.Server, timeout time.Duration, terminators ...TerminatorFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Run all terminations
	logger := slog.Default()
	logger.Info("performing pre-shutdown terminations")
	ctxTerm, cancelTerm := context.WithTimeout(context.Background(), timeout)
	var wg sync.WaitGroup
	wg.Add(len(terminators))
	for _, term := range terminators {
		go func(t TerminatorFunc) {
			defer wg.Done()
			if err := t(ctxTerm); err != nil {
				logger.Error(fmt.Sprintf("failed to terminate a resource, error: %v", err.Error()))
			}
		}(term)
	}

	wg.Wait()
	cancelTerm()

	// so server terminator doesn't have to be supplied
	if server != nil {
		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		if err := server.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("failed to shutdown http server, error: %v", err.Error()))
		}
	}

	logger.Info("shutting down")
	os.Exit(0)
}
