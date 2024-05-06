package kafka

import (
	"context"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

// NewProducer creates a new kafka producer
func NewProducer(client *kgo.Client, topic string) *Producer {
	return &Producer{client: client, topic: topic}
}

// SendMessage sends a message to the kafka topic and returns an error if any
func (p *Producer) SendMessage(ctx context.Context, message []byte) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var produceErr error
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: message}, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			produceErr = err
		}
	})
	wg.Wait()

	return produceErr
}

// Close closes the kafka producer
func (p *Producer) Close() {
	p.client.Close()
}
