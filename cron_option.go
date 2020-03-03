package dcron

import "time"

type CronOption func(c *Cron)

func WithKey(key string) CronOption {
	return func(c *Cron) {
		c.key = key
	}
}

func WithHostname(hostname string) CronOption {
	return func(c *Cron) {
		c.hostname = hostname
	}
}

func WithAtomic(atomic Atomic) CronOption {
	return func(c *Cron) {
		c.atomic = atomic
	}
}

func WithLocation(loc *time.Location) CronOption {
	return func(c *Cron) {
		c.location = loc
	}
}
