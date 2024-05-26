package utils

import (
	"fmt"

	cfgUtil "github.com/kartpop/cruncan/backend/pkg/config"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
	"github.com/kartpop/cruncan/backend/two/config"
	"github.com/wiremock/go-wiremock"
)

var EnvConfig = cfgUtil.LoadConfigOrPanic[config.Model]("../../config", "config")
var WiremockClient = wiremock.NewClient("http://localhost:8099")
var kafkaClient = initKafkaClient()
var OneRequestProducer = kafkaClient.NewProducer(EnvConfig.Kafka.OneRequestTopic.Name)

// TODO: client is not closed anywhere yet after initialization, refer one service tests
func initKafkaClient() *kafkaUtil.Client {
	kafkaClient, err := kafkaUtil.NewClient(EnvConfig.Kafka.Common)
	if err != nil {
		panic(fmt.Sprintf("failed to create kafka client: %v", err))
	}
	return kafkaClient
}
