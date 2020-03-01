package dcron

import (
	"context"
	"os"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	key      string
	hostname string
	cron     *cron.Cron
	mutex    Mutex
	jobs     []*innerJob
}

func NewCron(options ...CronOption) *Cron {
	ret := &Cron{
		cron: cron.New(cron.WithSeconds(), cron.WithLogger(cron.DiscardLogger)),
	}
	ret.hostname, _ = os.Hostname()
	for _, option := range options {
		option(ret)
	}
	return ret
}

func (c *Cron) AddJob(job Job) error {
	j := &innerJob{
		Job:  job,
		cron: c,
	}
	entryID, err := c.cron.AddJob(j.Spec(), j)
	if err != nil {
		return err
	}
	j.entryID = entryID
	c.jobs = append(c.jobs, j)
	return nil
}

func (c *Cron) Start() {
	c.cron.Start()
}

func (c *Cron) Stop() context.Context {
	return c.cron.Stop()
}

func (c *Cron) Run() {
	c.cron.Run()
}

func (c *Cron) Key() string {
	return c.key
}

func (c *Cron) Hostname() string {
	return c.hostname
}

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

func WithMutex(mutex Mutex) CronOption {
	return func(c *Cron) {
		c.mutex = mutex
	}
}
