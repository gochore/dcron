package dcron

import (
	"context"
	"time"
)

type GroupJob struct {
	cron *Cron
	key string
	spec string
	funcs []RunFunc
	funcKeys []string
}

func (g GroupJob) Key() string {
	panic("implement me")
}

func (g GroupJob) Spec() string {
	panic("implement me")
}

func (g GroupJob) Run(ctx context.Context) error {
	var planAt time.Time
	if task, ok := TaskFromContext(ctx); ok {
		planAt = task.PlanAt
	}

	for i, funcKey := range g.funcKeys {
		fn := g.funcs[i]
		key := genKey(planAt, g.cron.key, g.key, funcKey)
		if g.cron.setIfNotExists(key) {
			fn(ctx)
		}
	}

	panic("implement me")
}

func (g GroupJob) Options() []JobOption {
	panic("implement me")
}


