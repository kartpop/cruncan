package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type Client struct {
	client *kgo.Client
}

// NewClient creates a new kafka client
func NewClient(config *Config) (*Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.BootstrapServers...),
		kgo.ConsumerGroup(config.GroupId),
	)

	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

// Close closes the kafka client
func (c *Client) Close() {
	c.client.Close()
}

// NewProducer creates a new kafka producer
func (c *Client) NewProducer(topic string) *Producer {
	return NewProducer(c.client, topic)
}

// NewConsumer creates a new kafka consumer
func (c *Client) NewConsumer(topic string) *Consumer {
	return NewConsumer(c.client, topic)
}
