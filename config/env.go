package config

import (
	"github.com/spf13/viper"

	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	envConfig struct {
		configurable Configurable
	}
)

func newEnvConfig(configurable Configurable) (Configurable, error) {
	configKeys := []string{
		"PORT",

		"EVENT_STORAGE_HOST",
		"EVENT_STORAGE_PORT",
		"EVENT_STORAGE_USERNAME",
		"EVENT_STORAGE_PASSWORD",
		"EVENT_STORAGE_NAME",

		"EVENT_BUS_HOST",
		"EVENT_BUS_PORT",
		"EVENT_BUS_USERNAME",
		"EVENT_BUS_PASSWORD",

		"KAFKA_HOST",
		"KAFKA_PORT",
	}

	var err error
	for _, configKey := range configKeys {
		err = viper.BindEnv(configKey)
	}
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return envConfig{
		configurable: configurable,
	}, nil
}

func (instance envConfig) Get() (config Config, err error) {
	config, err = getConfig()
	if err != nil {
		return Config{}, helper.WrapError(err)
	}

	return config, nil
}
