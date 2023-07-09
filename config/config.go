package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	Config struct {
		Port                string
		EventStorage        EventStorage
		EventBus            EventBus
		Kafka               Kafka
		KSQLDB              KSQLDB
		KafkaSchemaRegistry KafkaSchemaRegistry
	}

	EventStorage struct {
		Host     string
		Port     string
		Username string
		Password string
		Name     string
	}

	EventBus struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	Kafka struct {
		Host string
		Port string
	}

	KSQLDB struct {
		Host string
		Port string
	}

	KafkaSchemaRegistry struct {
		Host string
		Port string
	}
)

const (
	missingConfigError = "config %s is missing"
)

var (
	instance *Config
)

func GetInstance() (Config, error) {
	configurable, err := newConfig()
	if err != nil {
		return Config{}, helper.WrapError(err)
	}

	if instance == nil {
		config, err := configurable.Get()
		if err != nil {
			return Config{}, helper.WrapError(err)
		}

		instance = &config
	}

	return *instance, nil
}

func newConfig() (Configurable, error) {
	defaultConfigurable := newDefaultConfig(nil)

	envConfigurable, err := newEnvConfig(defaultConfigurable)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	fileConfigurable, err := newFileConfig(envConfigurable)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return fileConfigurable, nil
}

func getConfig() (Config, error) {
	eventStoragePassword, err := getMandatoryString("EVENT_STORAGE_PASSWORD")
	if err != nil {
		return Config{}, helper.WrapError(err)
	}

	eventBusPassword, err := getMandatoryString("EVENT_BUS_PASSWORD")
	if err != nil {
		return Config{}, helper.WrapError(err)
	}

	return Config{
		Port: viper.GetString("PORT"),
		EventStorage: EventStorage{
			Host:     viper.GetString("EVENT_STORAGE_HOST"),
			Port:     viper.GetString("EVENT_STORAGE_PORT"),
			Username: viper.GetString("EVENT_STORAGE_USERNAME"),
			Password: eventStoragePassword,
			Name:     viper.GetString("EVENT_STORAGE_NAME"),
		},
		EventBus: EventBus{
			Host:     viper.GetString("EVENT_BUS_HOST"),
			Port:     viper.GetString("EVENT_BUS_PORT"),
			Username: viper.GetString("EVENT_BUS_USERNAME"),
			Password: eventBusPassword,
		},
		Kafka: Kafka{
			Host: viper.GetString("KAFKA_HOST"),
			Port: viper.GetString("KAFKA_PORT"),
		},
		KSQLDB: KSQLDB{
			Host: viper.GetString("KSQLDB_HOST"),
			Port: viper.GetString("KSQLDB_PORT"),
		},
		KafkaSchemaRegistry: KafkaSchemaRegistry{
			Host: viper.GetString("KSR_HOST"),
			Port: viper.GetString("KSR_PORT"),
		},
	}, nil
}

func getMandatoryString(key string) (string, error) {
	if !viper.IsSet(key) {
		return "", helper.WrapError(errors.New(fmt.Sprintf(missingConfigError, key)))
	}

	value := viper.GetString(key)
	return value, nil
}
