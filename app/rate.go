package app

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// MultiRateLimiter is a rate limiter which can handle cases where there
// are multiple rate limits that need to be enforced on a single resource
type MultiRateLimiter struct {
	limiters []*rate.Limiter
}

func NewMultiRateLimiter(limits []time.Duration, bursts []int) *MultiRateLimiter {
	limiters := make([]*rate.Limiter, len(limits))
	for i := range limits {
		limiters[i] = rate.NewLimiter(rate.Every(limits[i]), bursts[i])
	}
	return &MultiRateLimiter{limiters: limiters}
}

func (mrl *MultiRateLimiter) Wait(ctx context.Context) error {
	return mrl.WaitN(ctx, 1)
}

func (mrl *MultiRateLimiter) WaitN(ctx context.Context, n int) error {
	for i := range mrl.limiters {
		err := mrl.limiters[i].WaitN(ctx, n)
		if err != nil {
			return err
		}
	}
	return nil
}
