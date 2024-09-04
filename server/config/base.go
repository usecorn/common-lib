package config

import (
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

type BaseConfig struct {
	Name           string `env:"NAME" env-default:""`
	Version        string `env:"VERSION" env-default:""`
	Verbosity      string `env:"VERBOSITY" env-default:"INFO"`
	FluentDLogging bool   `env:"FLUENT_D_LOGGING" env-default:"false"`
	Production     bool   `env:"PRODUCTION" env-default:"false"`
}

func (c BaseConfig) GetLogger() *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(c.Verbosity)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(lvl)
	}

	logger.SetReportCaller(true)
	if c.FluentDLogging {
		logger.SetFormatter(joonix.NewFormatter())
	}
	return logger
}
