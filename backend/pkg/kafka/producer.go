package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/kartpop/cruncan/backend/pkg/util"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

// NewProducer creates a new kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		util.Fatal("failed to create kafka client for topic: %s, error: %v", topic, err)
	}

	return &Producer{client: client, topic: topic}
}

// SendMessage sends a message to the kafka topic
func (p *Producer) SendMessage(ctx context.Context, message []byte) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var produceErr error
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, func(_ *kgo.Record, err error) {
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
