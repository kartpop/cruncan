package config

import (
	gormUtil "github.com/kartpop/cruncan/backend/pkg/database/gorm"
)

type Model struct {
	Env                        string              `mapstructure:"ENVIRONMENT"`
	PodIP                      string              `mapstructure:"POD_IP"`
	Version                    string              `mapstructure:"VERSION"`
	LogLevel                   string              `mapstructure:"LOG_LEVEL"`
	Server   ServerConfig     `mapstructure:"SERVER"`
	Database *gormUtil.Config `mapstructure:"DATABASE"`
	Kafka    KafkaConfig      `mapstructure:"KAFKA"`
}

type ServerConfig struct {
	Addr         string `mapstructure:"ADDR"`
	WriteTimeout int    `mapstructure:"WRITE_TIMEOUT"`
	ReadTimeout  int    `mapstructure:"READ_TIMEOUT"`
	IdleTimeout  int    `mapstructure:"IDLE_TIMEOUT"`
}

type KafkaConfig struct {
	Common   KafkaCommon      `mapstructure:"COMMON"`
	OneTopic KafkaTopicConfig `mapstructure:"ONE_TOPIC"`
}

type KafkaCommon struct {
	BootstrapServers       string `mapstructure:"BOOTSTRAP_SERVERS"`
	SecurityProtocol       string `mapstructure:"SECURITY_PROTOCOL"`
	SslKeyLocation         string `mapstructure:"SSL_KEY_LOCATION"`
	SslCertificateLocation string `mapstructure:"SSL_CERTIFICATE_LOCATION"`
	GroupId                string `mapstructure:"GROUP_ID"`
	AutoOffsetReset        string `mapstructure:"AUTO_OFFSET_RESET"`
	LingerMs               string `mapstructure:"LINGER_MS"`
	BatchSize              string `mapstructure:"BATCH_SIZE"`
	LogLevel               string `mapstructure:"LOG_LEVEL"`
}

type KafkaTopicConfig struct {
	TopicName      string `mapstructure:"TOPIC_NAME"`
	PartitionCount int    `mapstructure:"PARTITION_COUNT"`
	ReplicaCount   int    `mapstructure:"REPLICA_COUNT"`
}
