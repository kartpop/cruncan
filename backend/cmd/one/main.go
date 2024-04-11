package main

import (
	"flag"

	"github.com/kartpop/cruncan/backend/config"
	cfgUtil "github.com/kartpop/cruncan/backend/pkg/config"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "address to listen for http traffic")
	dbServer := flag.String("dbserver", "", "database server name")
	dbPort := flag.Int("dbport", 0, "database server port")
	kafkaServers := flag.String("kafkaServers", "", "Kafka bootstrap servers")

	flag.Parse()

	var envConfig = cfgUtil.LoadConfigOrPanic[config.Model]()

	if httpAddr != "" {
		envConfig.OneHttpAddr = httpAddr
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
}
