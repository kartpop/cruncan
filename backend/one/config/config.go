package config

import (
	gormUtil "github.com/kartpop/cruncan/backend/pkg/database/gorm"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
)

type Model struct {
	Env      string           `mapstructure:"ENVIRONMENT"`
	PodIP    string           `mapstructure:"POD_IP"`
	Version  string           `mapstructure:"VERSION"`
	LogLevel string           `mapstructure:"LOG_LEVEL"`
	Server   ServerConfig     `mapstructure:"SERVER"`
	Database *gormUtil.Config `mapstructure:"DATABASE"`
	Kafka    KafkaConfig      `mapstructure:"KAFKA_CONFIG"`
}

type ServerConfig struct {
	Addr         string `mapstructure:"ADDR"`
	WriteTimeout int    `mapstructure:"WRITE_TIMEOUT"`
	ReadTimeout  int    `mapstructure:"READ_TIMEOUT"`
	IdleTimeout  int    `mapstructure:"IDLE_TIMEOUT"`
}

type KafkaConfig struct {
	Common kafkaUtil.Common `mapstructure:"COMMON"`
	OneRequestTopic kafkaUtil.Topic `mapstructure:"ONE_REQUEST_TOPIC"`
}

