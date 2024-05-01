package quickwit

import (
	"context"
	"math/rand"
	"time"
)

type BackoffConfig struct {
	// start backoff at this level
	MinBackoff time.Duration
	// increase exponentially to this level
	MaxBackoff time.Duration
	// give up after this many; zero means infinite retries
	MaxRetries int
}

// backoff implements exponential backoff with randomized wait times
type backoff struct {
	cfg          BackoffConfig
	ctx          context.Context
	numRetries   int
	nextDelayMin time.Duration
	nextDelayMax time.Duration
}

// newBackoff creates a backoff object. Pass a Context that can also terminate the operation.
func newBackoff(ctx context.Context, cfg BackoffConfig) *backoff {
	return &backoff{
		cfg:          cfg,
		ctx:          ctx,
		nextDelayMin: cfg.MinBackoff,
		nextDelayMax: doubleDuration(cfg.MinBackoff, cfg.MaxBackoff),
	}
}

// ongoing returns true if caller should keep going
func (b *backoff) ongoing() bool {
	// stop if Context has errored or max retry count is exceeded
	return b.ctx.Err() == nil && (b.cfg.MaxRetries == 0 || b.numRetries < b.cfg.MaxRetries)
}

// Wait sleeps for the backoff time then increases the retry count and backoff time
// Returns immediately if Context is terminated
func (b *backoff) wait() {
	// Increase the number of retries and get the next delay
	sleepTime := b.nextDelay()

	if b.ongoing() {
		select {
		case <-b.ctx.Done():
		case <-time.After(sleepTime):
		}
	}
}

func (b *backoff) nextDelay() time.Duration {
	b.numRetries++

	// Handle the edge case the min and max have the same value
	// (or due to some misconfig max is < min)
	if b.nextDelayMin >= b.nextDelayMax {
		return b.nextDelayMin
	}

	// Add a jitter within the next exponential backoff range
	sleepTime := b.nextDelayMin + time.Duration(rand.Int63n(int64(b.nextDelayMax-b.nextDelayMin)))

	// Apply the exponential backoff to calculate the next jitter
	// range, unless we've already reached the max
	if b.nextDelayMax < b.cfg.MaxBackoff {
		b.nextDelayMin = doubleDuration(b.nextDelayMin, b.cfg.MaxBackoff)
		b.nextDelayMax = doubleDuration(b.nextDelayMax, b.cfg.MaxBackoff)
	}

	return sleepTime
}

func doubleDuration(value time.Duration, max time.Duration) time.Duration {
	value = value * 2

	if value <= max {
		return value
	}

	return max
}
