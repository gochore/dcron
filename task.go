package dcron

import (
	"context"
	"time"
)

const (
	keyContextTask = "dcron/task"
)

type Task struct {
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

func TaskFromContext(ctx context.Context) (Task, bool) {
	if ctx == nil {
		return Task{}, false
	}
	task, ok := ctx.Value(keyContextTask).(Task)
	return task, ok
}
