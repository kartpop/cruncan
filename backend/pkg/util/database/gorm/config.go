package gorm

type Config struct {
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
