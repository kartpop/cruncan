package main

import (
	"errors"
	"flag"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
	"github.com/kartpop/cruncan/backend/config"
	cfgUtil "github.com/kartpop/cruncan/backend/pkg/config"
	"github.com/kartpop/cruncan/backend/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Application struct {
	name string
	cfg  *config.Model
	onePostHandler http.HandlerFunc
}

func NewApplication(name string, cfg *config.Model) *Application {
	return &Application{
		name: name,
		cfg:  cfg,
	}
}

func (app *Application) Run() []util.TerminatorFunc {
	return []util.TerminatorFunc{}
}

func (app *Application) routes() http.Handler {
	mux := flow.New()
	mux.HandleFunc("/one", app.onePostHandler, http.MethodPost)

	return mux
}

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "address to listen for http traffic")
	dbServer := flag.String("dbserver", "", "database server name")
	dbPort := flag.Int("dbport", 0, "database server port")
	kafkaServers := flag.String("kafkaServers", "", "Kafka bootstrap servers")

	flag.Parse()

	var envConfig = cfgUtil.LoadConfigOrPanic[config.Model]()

	if httpAddr != "" {
		envConfig.Server.HttpAddr = httpAddr
	}
	if *dbServer != "" {
		envConfig.Database.Server = *dbServer
	}
	if *dbPort != 0 {
		envConfig.Database.Port = *dbPort
	}
	if *kafkaServers != "" {
		envConfig.Kafka.Common.BootstrapServers = *kafkaServers
	}

	app := NewApplication("One", envConfig)

	server := &http.Server{
		Addr:         envConfig.Server.HttpAddr,
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
}
