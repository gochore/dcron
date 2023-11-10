package dcron

import (
	"context"
	"time"
)

// CronOption represents a modification to the default behavior of a Cron.
type CronOption func(c *Cron)

// WithKey overrides the key of the cron.
func WithKey(key string) CronOption {
	return func(c *Cron) {
		c.key = key
	}
}

// WithHostname overrides the hostname of the cron instance.
func WithHostname(hostname string) CronOption {
	return func(c *Cron) {
		c.hostname = hostname
	}
}

// WithAtomic uses the provided Atomic.
func WithAtomic(atomic Atomic) CronOption {
	return func(c *Cron) {
		c.atomic = atomic
	}
}

// WithLocation overrides the timezone of the cron instance.
func WithLocation(loc *time.Location) CronOption {
	return func(c *Cron) {
		c.location = loc
	}
}

// WithContext sets the root context of the cron instance.
// It will be used as the parent context of all tasks,
// and when the context is done, the cron will be stopped.
func WithContext(ctx context.Context) CronOption {
	return func(c *Cron) {
		c.context, c.contextCancel = context.WithCancel(ctx)
	}
}
