package dcron

import (
	"context"
	"time"
)

type Context struct {
	context.Context
	Key     string
	CronKey string
	JobKey  string
	PlanAt  time.Time
	BeginAt *time.Time
	EndAt   *time.Time
	Return  error
	Skipped bool
	Missed  bool
}
