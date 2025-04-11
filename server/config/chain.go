package config

import (
	"time"

	"golang.org/x/time/rate"
)

type Chain struct {
	MaxScanBlocks     uint64        `env:"MAX_SCAN_BLOCKS" env-default:"100"`
	RPCTimeout        time.Duration `env:"RPC_TIMEOUT" env-default:"5s"`
	RPCMaxConcurrency int64         `env:"RPC_MAX_CONCURRENCY" env-default:"5"`
	RPCRateLimit      int           `env:"RPC_RATE_LIMIT" env-default:"10"`
	RPCRateDur        time.Duration `env:"RPC_RATE_DUR" env-default:"1s"`
	RPCMaxRetries     int           `env:"RPC_MAX_RETRY" env-default:"5"`
	RPCRetryDelay     time.Duration `env:"RPC_RETRY_DELAY" env-default:"1s"`
	RPCURL            string        `env:"RPC_URL" env-default:""`
	CornRPCURL        string        `env:"CORN_RPC_URL" env-default:""`
	LagBlocks         uint64        `env:"LAG_BLOCKS" env-default:"9"`
}

func (c Chain) GetRateLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Limit(c.RPCRateLimit), c.RPCRateLimit)
}
