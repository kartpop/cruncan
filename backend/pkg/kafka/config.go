package kafka

type Config struct {
	Common Common `mapstructure:"COMMON"`
	Topic  Topic  `mapstructure:"TOPIC"`
}

type Common struct {
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

type Topic struct {
	Name           string `mapstructure:"NAME"`
	PartitionCount int    `mapstructure:"PARTITION_COUNT"`
	ReplicaCount   int    `mapstructure:"REPLICA_COUNT"`
}
