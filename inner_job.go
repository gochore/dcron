package dcron

import (
	"context"
	"fmt"
	"time"

	"github.com/gochore/pt"
	"github.com/robfig/cron/v3"
)

type innerJob struct {
	Job
	cron    *Cron
	entryID cron.EntryID
}

func (j *innerJob) Run() {
	c := j.cron
	entry := c.cron.Entry(j.entryID)
	planAt := entry.Prev
	nextAt := entry.Next
	key := fmt.Sprintf("dcron:%s.%s@%d", c.key, j.Key(), planAt.Unix())

	task := Task{
		Key:    key,
		Cron:   *c,
		Job:    j.Job,
		PlanAt: planAt,
	}
	task.Context, _ = context.WithDeadline(context.Background(), nextAt)

	skip := j.Before(task)
	if skip {
		task.Skipped = true
	}

	if !task.Skipped {
		if j.cron.mutex == nil || j.cron.mutex.SetIfNotExists(task.Key, c.hostname) {
			task.BeginAt = pt.Time(time.Now())
			if err := j.Job.Run(); err != nil {
				task.Return = err
			}
			task.EndAt = pt.Time(time.Now())
		} else {
			task.Missed = true
		}
	}

	j.After(task)
}

func (j *innerJob) Cron() *Cron {
	return j.cron
}
