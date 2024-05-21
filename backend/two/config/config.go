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
	Database *gormUtil.Config `mapstructure:"DATABASE"`
	Kafka    KafkaConfig      `mapstructure:"KAFKA_CONFIG"`
	Auth     AuthConfig       `mapstructure:"AUTH_CONFIG"`
	Three    ThreeConfig      `mapstructure:"THREE_CONFIG"`
}

type KafkaConfig struct {
	Common          *kafkaUtil.Config `mapstructure:"COMMON"`
	OneRequestTopic kafkaUtil.Topic   `mapstructure:"ONE_REQUEST_TOPIC"`
}

type AuthConfig struct {
	ClientID     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	TokenURL     string `mapstructure:"TOKEN_URL"`
}

type ThreeConfig struct {
	Url string `mapstructure:"URL"`
}
