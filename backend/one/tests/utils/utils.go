package utils

import (
	"fmt"

	"github.com/kartpop/cruncan/backend/one/config"
	cfgUtil "github.com/kartpop/cruncan/backend/pkg/config"
	gormUtil "github.com/kartpop/cruncan/backend/pkg/database/gorm"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
	"gorm.io/gorm"
)

var EnvConfig = cfgUtil.LoadConfigOrPanic[config.Model]("../../config", "config")

func InitKafkaClient() *kafkaUtil.Client {
	kafkaClient, err := kafkaUtil.NewClient(EnvConfig.Kafka.Common)
	if err != nil {
		panic(fmt.Sprintf("failed to create kafka client: %v", err))
	}

	return kafkaClient
}

func InitGorm() *gorm.DB {
	gormClient, err := gormUtil.NewGormClient(EnvConfig.Database)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize database connection: %v", err))
	}

	return gormClient
}
