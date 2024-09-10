package config

import (
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

type SentryConfig struct {
	Enabled         bool          `env:"SENTRY_ENABLED" env-default:"true"`
	DSN             string        `env:"SENTRY_DSN" env-default:""`
	EnableTrace     bool          `env:"SENTRY_ENABLE_TRACE" env-default:"true"`
	SampleRate      float64       `env:"SENTRY_SAMPLE_RATE" env-default:"1.0"` // Set to 1.0 to capture 100% of transactions
	Repanic         bool          `env:"SENTRY_REPANIC" env-default:"true"`
	Timeout         time.Duration `env:"SENTRY_TIMEOUT" env-default:"2s"`
	WaitForDelivery bool          `env:"SENTRY_WAIT_FOR_DELIVERY" env-default:"false"`
}

func (sc SentryConfig) SentryEnabled() bool {
	return len(sc.DSN) != 0 && sc.Enabled
}

func (sc SentryConfig) SetupSentry(name, version, env string) error {
	if !sc.SentryEnabled() {
		return nil
	}
	return sentry.Init(sentry.ClientOptions{
		Dsn:              sc.DSN,
		EnableTracing:    sc.EnableTrace,
		TracesSampleRate: sc.SampleRate,
		ServerName:       name,
		Release:          version,
		Environment:      env,
	})
}

func (sc SentryConfig) NewGinMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic:         sc.Repanic,
		WaitForDelivery: sc.WaitForDelivery,
		Timeout:         sc.Timeout,
	})
}

func NewSentryConfig() (SentryConfig, error) {
	sentryConfig := SentryConfig{}
	err := cleanenv.ReadEnv(&sentryConfig)
	return sentryConfig, err
}
