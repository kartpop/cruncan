package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfigOrPanic[T any](pathParams ...string) *T {
	var filePath, fileName string
	if pathParams == nil {
		// this is a default setup for prod usage
		filePath, fileName = "config", "config"
	} else if len(pathParams) == 2 {
		// this is very useful to use in tests when you want to specify which config file in which location to use
		filePath, fileName = pathParams[0], pathParams[1]
	} else {
		panic(fmt.Errorf("incorrect amount of passed params: %v", pathParams))
	}

	loadConfig, err := loadConfig[T](filePath, fileName)
	if err != nil {
		panic(fmt.Errorf("failed to load config, error: %v", err.Error()))
	}

	return loadConfig
}

func loadConfig[T any](filePath, fileName string) (*T, error) {
	err := initViper(filePath, fileName)
	if err != nil {
		return nil, err
	}

	return loadConfigObject[T]()
}

func initViper(filePath, fileName string) error {
	viper.SetConfigName(fileName)
	viper.AddConfigPath(filePath)
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

func loadConfigObject[T any]() (*T, error) {
	var config *T
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
