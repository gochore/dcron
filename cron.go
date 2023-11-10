package dcron

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// CronMeta is a read only wrapper for Cron.
type CronMeta interface {
	// Key returns the unique key of the cron.
	Key() string
	// Hostname returns current hostname.
	Hostname() string
	// Statistics returns statistics info of the cron's all jobs.
	Statistics() Statistics
	// Jobs returns the cron's all jobs as JobMeta.
	Jobs() []JobMeta
}

// Cron keeps track of any number of jobs, invoking the associated func as specified.
type Cron struct {
	key           string
	hostname      string
	cron          *cron.Cron
	atomic        Atomic
	jobs          []*innerJob
	location      *time.Location
	context       context.Context
	contextCancel context.CancelFunc
}

// NewCron returns a cron with specified options.
func NewCron(options ...CronOption) *Cron {
	ret := &Cron{
		location: time.Local,
	}
	ret.hostname, _ = os.Hostname()
	for _, option := range options {
		option(ret)
	}

	ret.cron = cron.New(
		cron.WithSeconds(),
		cron.WithLogger(cron.DiscardLogger),
		cron.WithLocation(ret.location),
	)
	return ret
}

// AddJobs helps to add multiple jobs.
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
	if job.Key() == "" {
		return errors.New("empty key")
	}

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

// Start the cron scheduler in its own goroutine, or no-op if already started.
func (c *Cron) Start() {
	if c.context != nil {
		go func() {
			<-c.context.Done()
			c.Stop()
		}()
	}
	c.cron.Start()
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
// A context is returned so the caller can wait for running jobs to complete.
func (c *Cron) Stop() context.Context {
	if c.contextCancel != nil {
		c.contextCancel()
	}
	return c.cron.Stop()
}

// Run the cron scheduler, or no-op if already running.
func (c *Cron) Run() {
	if c.context != nil {
		go func() {
			<-c.context.Done()
			c.Stop()
		}()
	}
	c.cron.Run()
}

// Key implements CronMeta.Key
func (c *Cron) Key() string {
	return c.key
}

// Hostname implements CronMeta.Hostname
func (c *Cron) Hostname() string {
	return c.hostname
}

// Statistics implements CronMeta.Statistics
func (c *Cron) Statistics() Statistics {
	ret := Statistics{}
	for _, j := range c.jobs {
		ret = ret.Add(j.statistics)
	}
	return ret
}

// Jobs implements CronMeta.Jobs
func (c *Cron) Jobs() []JobMeta {
	var ret []JobMeta
	for _, j := range c.jobs {
		ret = append(ret, j)
	}
	return ret
}
