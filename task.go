package dcron

import (
	"context"
	"time"
)

type Task struct {
	context.Context
	Key     string
	Cron    Cron
	Job     Job
	PlanAt  time.Time
	BeginAt *time.Time
	EndAt   *time.Time
	Return  error
	Skipped bool
	Missed  bool
}
