package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"github.com/kartpop/cruncan/backend/reference/retry"
)

func main() {
	ctx := context.Background()

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", "localhost", "6379"),
				redis.DialPassword(""),
				redis.DialUseTLS(false),
			)
			if err != nil {
				slog.Default().Error(fmt.Sprintf("failed to dial redis, error: %v", err))
			}
			return c, err
		},
	}

	redisSync := redsync.New(redigo.NewPool(pool))

	retryJob := retry.NewJob(slog.Default(), redisSync)
	retryJob.Start(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("stopping retry job...")
	retryJob.Stop(ctx)
}
