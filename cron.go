package dcron

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/robfig/cron/v3"
)

type CronMeta interface {
	Key() string
	Hostname() string
}

type Cron struct {
	key      string
	hostname string
	cron     *cron.Cron
	atomic   Atomic
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

func (c *Cron) AddJobs(jobs ...Job) error {
	var errs []string
	for _, job := range jobs {
		if err := c.addJob(job); err != nil {
			errs = append(errs, fmt.Sprintf("add job %s: %v", job.Key(), err))
		}
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func (c *Cron) addJob(job Job) error {
	for _, j := range c.jobs {
		if j.key == job.Key() {
			return errors.New("added already")
		}
	}

	j := &innerJob{
		cron:        c,
		entryGetter: c.cron,
		key:         job.Key(),
		spec:        job.Spec(),
		run:         job.Run,
	}

	for _, option := range job.Options() {
		option(j)
	}
	if j.retryTimes < 1 {
		j.retryTimes = 1
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
