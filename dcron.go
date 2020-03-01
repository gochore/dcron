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
	cron     cron.Cron
	mutex    Mutex
	jobs     []*addedJob
}

func (c *Cron) AddJob(job Job) error {
	j := &addedJob{
		Job: job,
	}
	entryID, err := c.cron.AddJob(j.Spec(), j)
	if err != nil {
		return err
	}
	j.entryID = entryID
	c.jobs = append(c.jobs, j)
	return nil
}

type addedJob struct {
	Job
	cron    *Cron
	entryID cron.EntryID
}

func (j *addedJob) Run() {
	c := j.cron
	entry := c.cron.Entry(j.entryID)
	key := fmt.Sprintf("dcron:%s.%s@%d", c.key, j.Key(), entry.Next.Unix())

	ctx := JobContext{
		Context: context.TODO(),
		Key:     key,
		CronKey: c.key,
		JobKey:  j.Key(),
		PlanAt:  entry.Next,
	}

	skip := j.Before(ctx)
	if skip {
		ctx.Skipped = true
	}

	if !ctx.Skipped {
		if j.cron.mutex.SetIfNotExists(key, c.hostname) {
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
