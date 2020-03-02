package dcron

import (
	"context"
	"time"
)

type Task struct {
	ctx        context.Context
	Key        string
	Cron       CronMeta
	Job        JobMeta
	PlanAt     time.Time
	BeginAt    *time.Time
	EndAt      *time.Time
	Return     error
	Skipped    bool
	Missed     bool
	TriedTimes int
}

func (t Task) Deadline() (deadline time.Time, ok bool) {
	return t.ctx.Deadline()
}

func (t Task) Done() <-chan struct{} {
	return t.ctx.Done()
}

func (t Task) Err() error {
	return t.ctx.Err()
}

func (t Task) Value(key interface{}) interface{} {
	return t.ctx.Value(key)
}
