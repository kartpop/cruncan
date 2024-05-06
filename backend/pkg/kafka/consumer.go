package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/exp/slog"
)

type ConsumerHandler interface {
	Handle([]byte, string) error
}

type Consumer struct {
	client *kgo.Client
	topic  string
}

// NewConsumer creates a new kafka consumer
func NewConsumer(client *kgo.Client, topic string) *Consumer {
	client.AddConsumeTopics(topic)
	return &Consumer{client: client, topic: topic}
}

// Start starts the kafka consumer and calls the handler for each message for the topic until the context is cancelled.
// The poll loop is inside a goroutine, so it will not block the caller.
func (c *Consumer) Start(ctx context.Context, handler ConsumerHandler) {
	go func() {
		for {
			fetches := c.client.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				for _, err := range errs {
					if errors.Is(err.Err, kgo.ErrClientClosed) {
						return
					}
				}
				// All errors are retried internally when fetching, but non-retriable errors are
				// returned from polls so that users can notice and take action.
				slog.ErrorContext(ctx, fmt.Sprint(errs))
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				if record.Topic == c.topic {
					// ignore error here, handle it in the handler
					_ = handler.Handle(record.Value, record.Topic)
				}
			}
		}
	}()
}