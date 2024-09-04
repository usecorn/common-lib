package config

type BaseConfig struct {
	Name           string `env:"NAME" env-default:""`
	Version        string `env:"VERSION" env-default:""`
	Verbosity      string `env:"VERBOSITY" env-default:"INFO"`
	FluentDLogging bool   `env:"FLUENT_D_LOGGING" env-default:"false"`
}
