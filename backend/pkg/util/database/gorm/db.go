package gorm

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewGormClient(config *Config) (*gorm.DB, error) {
	conn, err := createConnection(config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connect(dialector gorm.Dialector, config *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: NewSlogLogger(config.LogLevel),
	})
	if err != nil {
		return nil, err
	}

	err = db.Use(tracing.NewPlugin())

	err = setTimeOutsAndMaxConns(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setTimeOutsAndMaxConns(db *gorm.DB) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(20)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDb.SetMaxOpenConns(20)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(time.Duration(30) * time.Second)

	return nil
}

func createConnection(config *Config) (*gorm.DB, error) {
	dbDsn := makeDsn(config)
	pgDialect := postgres.Open(dbDsn)
	return connect(pgDialect, config)
}

func makeDsn(config *Config) string {
	dsnTemplate := "host=%v user=%v password=%v dbname=%v port=%d sslmode=%v TimeZone=UTC"

	return fmt.Sprintf(
		dsnTemplate,
		config.Server,
		config.User,
		config.Password,
		config.Name,
		config.Port,
		string(config.SslMode),
	)
}
