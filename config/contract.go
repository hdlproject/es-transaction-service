package config

type (
	Configurable interface {
		Get() (Config, error)
	}
)
