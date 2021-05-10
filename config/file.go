package config

import (
	"errors"
	"github.com/hdlproject/es-transaction-service/helper"
	"github.com/spf13/viper"
)

type (
	fileConfig struct {
		configurable Configurable
	}
)

func newFileConfig(configurable Configurable) (Configurable, error) {
	configFilepath := viper.GetString("CONFIG_FILEPATH")
	configFilename := viper.GetString("CONFIG_FILENAME")

	viper.SetConfigName(configFilename)
	viper.SetConfigType("env")
	viper.AddConfigPath(configFilepath)

	err := viper.ReadInConfig()
	var errConfigFileNotFound *viper.ConfigFileNotFoundError
	if err != nil && errors.As(err, &errConfigFileNotFound) {
		return nil, helper.WrapError(err)
	}

	return fileConfig{
		configurable: configurable,
	}, nil
}

func (instance fileConfig) Get() (config Config, err error) {
	config, err = getConfig()
	if err != nil {
		return Config{}, helper.WrapError(err)
	}

	return config, nil
}
