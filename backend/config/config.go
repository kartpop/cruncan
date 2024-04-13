package config

type Model struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	HttpAddr     string `mapstructure:"HTTP_ADDR"`
	WriteTimeout int    `mapstructure:"WRITE_TIMEOUT"`
	ReadTimeout  int    `mapstructure:"READ_TIMEOUT"`
	IdleTimeout  int    `mapstructure:"IDLE_TIMEOUT"`
}

type DatabaseConfig struct {
	Server   string  `mapstructure:"DATABASE_SERVER"`
	Port     int     `mapstructure:"DATABASE_PORT"`
	Name     string  `mapstructure:"DATABASE_NAME"`
	User     string  `mapstructure:"DATABASE_USER"`
	Password string  `mapstructure:"DATABASE_PASSWORD"`
	LogLevel string  `mapstructure:"DATABASE_LOG_LEVEL"`
	SslMode  SslMode `mapstructure:"DATABASE_SSL_MODE"`
}

// https://www.postgresql.org/docs/current/libpq-ssl.html
type SslMode string

const (
	SslModeDisable    SslMode = SslMode("disable")
	SslModeAllow      SslMode = SslMode("allow")
	SslModePrefer     SslMode = SslMode("prefer")
	SslModeRequire    SslMode = SslMode("require")
	SslModeVerifyCa   SslMode = SslMode("verify-ca")
	SslModeVerifyFull SslMode = SslMode("verify-full")
)

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
