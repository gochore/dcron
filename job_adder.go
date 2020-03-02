package dcron

import (
	"errors"
	"fmt"
	"strings"
)

type JobAdder struct {
	errs []error
	cron *Cron
}

func NewJobAdder(cron *Cron) *JobAdder {
	return &JobAdder{
		cron: cron,
	}
}

func (adder *JobAdder) Add(job Job) {
	if adder.cron == nil {
		adder.errs = append(adder.errs, fmt.Errorf("add job %s: cron is nil", job.Key()))
		return
	}
	if err := adder.cron.AddJob(job); err != nil {
		adder.errs = append(adder.errs, fmt.Errorf("add job %s: %w", job.Key(), err))
	}
}

func (adder *JobAdder) Err() error {
	if len(adder.errs) == 0 {
		return nil
	}
	if len(adder.errs) == 1 {
		return adder.errs[0]
	}
	var errs []string
	for _, err := range adder.errs {
		errs = append(errs, err.Error())
	}
	return errors.New(strings.Join(errs, "; "))
}
