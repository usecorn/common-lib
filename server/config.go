package server

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" env-default:"5s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"15s"`
	Listen          string        `env:"LISTEN" env-default:":8080"`
	HealthCheckPath string        `env:"HEALTH_CHECK_PATH" env-default:"/health"`
}

func LoadServerConfigFromEnv() (ServerConfig, error) {
	conf := ServerConfig{}
	err := cleanenv.ReadEnv(&conf)

	return conf, err
}
