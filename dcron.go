package dcron

import (
	"context"
	"fmt"
	"time"

	"github.com/gochore/pt"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	key      string
	hostname string
	cron     *cron.Cron
	mutex    Mutex
	jobs     []*addedJob
}

func NewCron() *Cron {
	return &Cron{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (c *Cron) AddJob(job Job) error {
	j := &addedJob{
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

type addedJob struct {
	Job
	cron    *Cron
	entryID cron.EntryID
}

func (j *addedJob) Run() {
	c := j.cron
	planAt := c.cron.Entry(j.entryID).Prev
	key := fmt.Sprintf("dcron:%s.%s@%d", c.key, j.Key(), planAt.Unix())

	ctx := Context{
		Context: context.TODO(),
		Key:     key,
		CronKey: c.key,
		JobKey:  j.Key(),
		PlanAt:  planAt,
	}

	skip := j.Before(ctx)
	if skip {
		ctx.Skipped = true
	}

	if !ctx.Skipped {
		if j.cron.mutex == nil || j.cron.mutex.SetIfNotExists(ctx.Key, c.hostname) {
			ctx.BeginAt = pt.Time(time.Now())
			if err := j.Job.Run(); err != nil {
				ctx.Return = err
			}
			ctx.EndAt = pt.Time(time.Now())
		} else {
			ctx.Missed = true
		}
	}

	j.After(ctx)
}
