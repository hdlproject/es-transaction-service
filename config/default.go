package config

import (
	"github.com/spf13/viper"

	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	defaultConfig struct {
		configurable Configurable
	}
)

func newDefaultConfig(configurable Configurable) Configurable {
	viper.SetDefault("PORT", "7777")

	viper.SetDefault("CONFIG_FILEPATH", ".")
	viper.SetDefault("CONFIG_FILENAME", ".env")

	viper.SetDefault("EVENT_STORE_HOST", "127.0.0.1")
	viper.SetDefault("EVENT_STORE_PORT", "27017")
	viper.SetDefault("EVENT_STORE_USERNAME", "root")
	viper.SetDefault("EVENT_STORE_NAME", "es-transaction-service")

	viper.SetDefault("EVENT_BUS_HOST", "127.0.0.1")
	viper.SetDefault("EVENT_BUS_PORT", "5672")
	viper.SetDefault("EVENT_BUS_USERNAME", "root")

	viper.SetDefault("KAFKA_HOST", "127.0.0.1")
	viper.SetDefault("KAFKA_PORT", "29092")

	return defaultConfig{
		configurable: configurable,
	}
}

func (instance defaultConfig) Get() (config Config, err error) {
	config, err = getConfig()
	if err != nil {
		return Config{}, helper.WrapError(err)
	}

	return config, nil
}
