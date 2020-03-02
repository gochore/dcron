package dcron

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
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

	ctx, cancel := context.WithDeadline(context.WithValue(context.Background(), keyContextTask, task), nextAt)
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
				task.Return = safeRun(ctx, j.run)
				task.TriedTimes++
				if task.Return == nil {
					break
				}
				if ctx.Err() != nil {
					break
				}
				if j.retryInterval != nil {
					interval := j.retryInterval(task.TriedTimes)
					deadline, _ := ctx.Deadline()
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

func safeRun(ctx context.Context, run RunFunc) (err error) {
	defer func() {
		if r := recover(); r != nil {
			pc := make([]uintptr, 16)
			n := runtime.Callers(0, pc)
			for _, p := range pc[:n] {
				fn := runtime.FuncForPC(p)
				if fn != nil {
					file, line := fn.FileLine(p)
					if !strings.Contains(fn.Name(), "runtime") {
						err = fmt.Errorf("panic(%v) at %s:%d", r, file, line)
						return
					}
				}
			}
			err = fmt.Errorf("panic(%v): %s", r, debug.Stack())
		}
	}()
	return run(ctx)
}
