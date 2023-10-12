package dcron

import (
	"context"
	"time"
)

type ctxKey string

const (
	keyContextTask ctxKey = "dcron/task"
)

// Task is an execute of a job.
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

// TaskFromContext extracts a Task from a context,
// it is useful inner the Run function.
func TaskFromContext(ctx context.Context) (Task, bool) {
	if ctx == nil {
		return Task{}, false
	}
	task, ok := ctx.Value(keyContextTask).(Task)
	return task, ok
}
