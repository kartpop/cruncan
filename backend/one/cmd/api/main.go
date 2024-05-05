package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
	"github.com/kartpop/cruncan/backend/one/config"
	onerequest "github.com/kartpop/cruncan/backend/one/database/one_request"
	oneHttp "github.com/kartpop/cruncan/backend/one/http"
	cfgUtil "github.com/kartpop/cruncan/backend/pkg/config"
	gormUtil "github.com/kartpop/cruncan/backend/pkg/database/gorm"
	"github.com/kartpop/cruncan/backend/pkg/id"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
	"github.com/kartpop/cruncan/backend/pkg/otel"
	"github.com/kartpop/cruncan/backend/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var tracerName = "github.com/kartpop/cruncan/backend/one/cmd/api-tracer"
var meterName = "github.com/kartpop/cruncan/backend/one/cmd/api-meter"

type Application struct {
	name        string
	cfg         *config.Model
	oneHandler  *oneHttp.Handler
	kafkaClient *kafkaUtil.Client
}

func NewApplication(ctx context.Context, name string, cfg *config.Model) *Application {

	gormClient, err := gormUtil.NewGormClient(cfg.Database)
	if err != nil {
		util.Fatal("database not available on startup: %v", err)
	}

	idService, err := id.NewServiceFromIP(cfg.PodIP)
	if err != nil {
		util.Fatal("failed to create id service: %v", err)
	}

	kafkaClient, err := kafkaUtil.NewClient(cfg.Kafka.Common)
	if err != nil {
		util.Fatal("failed to create kafka client: %v", err)
	}
	oneRequestProducer := kafkaClient.NewProducer(cfg.Kafka.OneRequestTopic.Name)
	oneRequestRepo := onerequest.NewRepository(gormClient).WithTracing()
	oneHandler := oneHttp.NewHandler(ctx, oneRequestRepo, idService, oneRequestProducer)

	return &Application{
		name:       name,
		cfg:        cfg,
		oneHandler: oneHandler,
		kafkaClient: kafkaClient,
	}
}

func (app *Application) Run() []util.TerminatorFunc {
	return []util.TerminatorFunc{
		func(ctx context.Context) error {
			app.kafkaClient.Close()
			return nil
		},
	}
}

func (app *Application) routes() http.Handler {
	mux := flow.New()
	mux.HandleFunc("/one", app.oneHandler.Post, http.MethodPost)

	return mux
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

	if httpAddr != "" {
		envConfig.Server.Addr = httpAddr
	}
	if *dbServer != "" {
		envConfig.Database.Server = *dbServer
	}
	if *dbPort != 0 {
		envConfig.Database.Port = *dbPort
	}
	if *kafkaServers != "" {
		envConfig.Kafka.Common.BootstrapServers = []string{*kafkaServers}
	}

	app := NewApplication(ctx, "One", envConfig)

	server := &http.Server{
		Addr:         envConfig.Server.Addr,
		WriteTimeout: time.Duration(envConfig.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(envConfig.Server.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(envConfig.Server.IdleTimeout) * time.Second,
		Handler:      otelhttp.NewHandler(app.routes(), "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			util.Fatal(err.Error())
		}
	}()

	terminatorFunctions := app.Run()

	slog.InfoContext(ctx, fmt.Sprintf("%v is running!", app.name))

	util.GracefulShutdown(
		server, time.Second*5,
		terminatorFunctions...)
}
