package dcron

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

// CronMeta is a read only wrapper for innerJob.
type JobMeta interface {
	// Key returns the unique key of the job.
	Key() string
	// Key returns the spec of the job.
	Spec() string
	// Statistics returns statistics info of the job.
	Statistics() Statistics
}

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
	noMutex       bool
	statistics    Statistics
}

// Key implements JobMeta.Key.
func (j *innerJob) Key() string {
	return j.key
}

// Spec implements JobMeta.Spec.
func (j *innerJob) Spec() string {
	return j.spec
}

// Statistics implements JobMeta.Statistics.
func (j *innerJob) Statistics() Statistics {
	return j.statistics
}

func (j *innerJob) Run() {
	c := j.cron
	entry := j.entryGetter.Entry(j.entryID)
	planAt := entry.Prev
	nextAt := entry.Next
	key := genKey(planAt, c.key, j.key)

	task := Task{
		Key:        key,
		Cron:       c,
		Job:        j,
		PlanAt:     planAt,
		TriedTimes: 0,
	}
	atomic.AddInt64(&j.statistics.TotalTask, 1)

	ctx, cancel := context.WithDeadline(context.WithValue(context.Background(), keyContextTask, task), nextAt)
	defer cancel()

	if j.before != nil && j.before(task) {
		task.Skipped = true
		atomic.AddInt64(&j.statistics.SkippedTask, 1)
	}

	if !task.Skipped {
		if j.noMutex ||  j.cron.setIfNotExists(task.Key) {
			beginAt := time.Now()
			task.BeginAt = &beginAt

			for i := 0; i < j.retryTimes; i++ {
				task.Return = safeRun(ctx, j.run)
				atomic.AddInt64(&j.statistics.TotalRun, 1)
				if i > 0 {
					atomic.AddInt64(&j.statistics.RetriedRun, 1)
				}
				task.TriedTimes++
				if task.Return == nil {
					atomic.AddInt64(&j.statistics.PassedRun, 1)
					break
				}
				atomic.AddInt64(&j.statistics.FailedRun, 1)
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
			atomic.AddInt64(&j.statistics.MissedTask, 1)
		}
	}

	if j.after != nil {
		j.after(task)
	}

	if !task.Skipped && !task.Missed {
		if task.Return == nil {
			atomic.AddInt64(&j.statistics.PassedTask, 1)
		} else {
			atomic.AddInt64(&j.statistics.FailedTask, 1)
		}
	}
}

func safeRun(ctx context.Context, run RunFunc) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v: %s", r, debug.Stack())
		}
	}()
	return run(ctx)
}

func genKey(t time.Time, keys ...string) string {
	return fmt.Sprintf("dcron:%s@%d", strings.Join(keys, "."), t.Unix())
}