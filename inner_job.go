package dcron

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type innerJob struct {
	cron          *Cron
	entryID       cron.EntryID
	entryGetter   entryGetter
	key           string
	spec          string
	before        BeforeFunc
	run           RunFunc
	after         AfterFunc
	retryTimes    int
	retryInterval RetryInterval
}

func (j *innerJob) Key() string {
	return j.key
}

func (j *innerJob) Spec() string {
	return j.spec
}

func (j *innerJob) Run() {
	c := j.cron
	entry := j.entryGetter.Entry(j.entryID)
	planAt := entry.Prev
	nextAt := entry.Next
	key := fmt.Sprintf("dcron:%s.%s@%d", c.key, j.key, planAt.Unix())

	task := Task{
		Key:        key,
		Cron:       c,
		Job:        j,
		PlanAt:     planAt,
		TriedTimes: 0,
	}

	var cancel context.CancelFunc
	task.ctx, cancel = context.WithDeadline(context.Background(), nextAt)
	defer cancel()

	skip := false
	if j.before != nil && j.before(task) {
		skip = true
	}

	if skip {
		task.Skipped = true
	}

	if !task.Skipped {
		if j.cron.mutex == nil || j.cron.mutex.SetIfNotExists(task.Key, c.hostname) {
			beginAt := time.Now()
			task.BeginAt = &beginAt

			for i := 0; i < j.retryTimes; i++ {
				task.Return = j.run(task)
				task.TriedTimes++
				if task.Return == nil {
					break
				}
				if task.Err() != nil {
					break
				}
				if j.retryInterval != nil {
					interval := j.retryInterval(task.TriedTimes)
					deadline, _ := task.Deadline()
					if -time.Since(deadline) < interval {
						break
					}
					time.Sleep(interval)
				}
			}

			endAt := time.Now()
			task.EndAt = &endAt
		} else {
			task.Missed = true
		}
	}

	if j.after != nil {
		j.after(task)
	}
}

func (j *innerJob) Cron() *Cron {
	return j.cron
}
