package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"time"

	cfgUtil "github.com/kartpop/cruncan/backend/pkg/config"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
	"github.com/kartpop/cruncan/backend/pkg/otel"
	"github.com/kartpop/cruncan/backend/pkg/util"
	"github.com/kartpop/cruncan/backend/two/config"
	"github.com/kartpop/cruncan/backend/two/onerequest"
)

var tracerName = "github.com/kartpop/cruncan/backend/two/cmd/consumer-tracer"
var meterName = "github.com/kartpop/cruncan/backend/two/cmd/consumer-meter"

type Application struct {
	ctx context.Context
	name string
	cfg  *config.Model
	kafkaClient *kafkaUtil.Client
	oneRequestConsumer *kafkaUtil.Consumer
	oneRequestKafkaHandler *onerequest.KafkaHandler
}

func NewApplication(ctx context.Context, name string, cfg *config.Model) *Application {
	kafkaClient, err := kafkaUtil.NewClient(cfg.Kafka.Common)
	if err != nil {
		util.Fatal("failed to create kafka client: %v", err)
	}

	oneRequestConsumer := kafkaClient.NewConsumer(cfg.Kafka.OneRequestTopic.Name)
	oneRequestKafkaHandler := onerequest.NewKafkaHandler()

	return &Application{
		ctx: ctx,
		name: name,
		cfg: cfg,
		kafkaClient: kafkaClient,
		oneRequestConsumer: oneRequestConsumer,
		oneRequestKafkaHandler: oneRequestKafkaHandler,
	}
}

func (app *Application) Run() []util.TerminatorFunc {
	app.oneRequestConsumer.Start(app.ctx, app.oneRequestKafkaHandler)

	return []util.TerminatorFunc{
		func(ctx context.Context) error {
			app.kafkaClient.Close()
			return nil
		},
	}
}

func main() {
	ctx, cancel := otel.Setup(tracerName, meterName)
	defer cancel()

	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "address to listen for http traffic")
	dbServer := flag.String("dbserver", "", "database server name")
	dbPort := flag.Int("dbport", 0, "database server port")
	kafkaServers := flag.String("kafkaServers", "", "Kafka bootstrap servers")

	flag.Parse()

	var envConfig = cfgUtil.LoadConfigOrPanic[config.Model]()

	if *dbServer != "" {
		envConfig.Database.Server = *dbServer
	}
	if *dbPort != 0 {
		envConfig.Database.Port = *dbPort
	}
	if *kafkaServers != "" {
		envConfig.Kafka.Common.BootstrapServers = []string{*kafkaServers}
	}

	app := NewApplication(ctx, "two-consumer", envConfig)

	terminatorFunctions := app.Run()

	slog.InfoContext(ctx, fmt.Sprintf("%v is running!", app.name))

	util.GracefulShutdown(
		nil, time.Second*5,
		terminatorFunctions...
	)
}