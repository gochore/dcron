package dcron

import "time"

// CronOption represents a modification to the default behavior of a Cron.
type CronOption func(c *Cron)

// WithLocation overrides the key of the cron.
func WithKey(key string) CronOption {
	return func(c *Cron) {
		c.key = key
	}
}

// WithLocation overrides the hostname of the cron instance.
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
