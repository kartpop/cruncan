package retry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type Job struct {
	logger *slog.Logger
	rs     *redsync.Redsync
	ticker *time.Ticker
}

func NewJob(logger *slog.Logger, rs *redsync.Redsync) *Job {
	return &Job{
		logger: logger,
		rs:     rs,
	}
}

func (j *Job) Start(ctx context.Context) {
	if j.ticker != nil {
		j.logger.WarnContext(ctx, "trying to start a transaction retry job that is already running")
		return
	}

	interval := 5 * time.Minute
	j.ticker = time.NewTicker(interval)

	go func() {
		for t := range j.ticker.C {
			j.logger.DebugContext(ctx, fmt.Sprintf("starting transaction retry job at %v", t))
			failedTxns := GetFailedTransactions()

			for _, failedTxn := range failedTxns {
				j.ProcessSingleTxn(ctx, failedTxn)
			}

			j.logger.InfoContext(ctx, "finished transaction retry job")
		}
	}()
}

func (j *Job) ProcessSingleTxn(ctx context.Context, failedTxn Transaction) {
	mutex := j.rs.NewMutex(failedTxn.ID, redsync.WithExpiry(60*time.Second))
	if err := mutex.Lock(); err != nil {
		j.logger.ErrorContext(ctx, fmt.Sprintf("failed to obtain redsync mutex for transacation %s: %v", failedTxn.ID, err.Error()))
		return
	}

	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			j.logger.ErrorContext(ctx, fmt.Sprintf("failed to release redsync mutex for transaction %s: %v", failedTxn.ID, err.Error()))
		}
	}()

	// lock acquired, recheck transaction state
	if IsTransactionFailed(failedTxn) {
		j.logger.DebugContext(ctx, fmt.Sprintf("retrying transaction %s", failedTxn.ID))

		// business logic for retry
		// ...
	} else {
		j.logger.DebugContext(ctx, fmt.Sprintf("transaction %s already processed", failedTxn.ID))
	}
}

func (j *Job) Stop(ctx context.Context) error {
	if j.ticker != nil {
		j.ticker.Stop()
	}

	j.ticker = nil
	return nil
}
