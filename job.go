package dcron

import (
	"fmt"
)

type Job interface {
	Key() string
	Spec() string
	Before(task Task) (skip bool)
	Run() error
	After(task Task)
}

type BeforeFunc func(task Task) (skip bool)

type RunFunc func() error

type AfterFunc func(task Task)

type wrapJob struct {
	key    string
	spec   string
	before BeforeFunc
	run    RunFunc
	after  AfterFunc
}

func NewJob(key, spec string, run RunFunc, options ...JobOption) (Job, error) {
	if key == "" {
		return nil, fmt.Errorf("empty key")
	}
	if run == nil {
		return nil, fmt.Errorf("nil run")
	}
	ret := &wrapJob{
		key:  key,
		spec: spec,
		run:  run,
	}
	for _, option := range options {
		option(ret)
	}
	return ret, nil
}

func (j *wrapJob) Key() string {
	return j.key
}

func (j *wrapJob) Spec() string {
	return j.spec
}

func (j *wrapJob) Before(task Task) (skip bool) {
	if j.before != nil {
		return j.before(task)
	}
	return false
}

func (j *wrapJob) Run() error {
	return j.run()
}

func (j *wrapJob) After(task Task) {
	if j.after != nil {
		j.after(task)
	}
}

type JobOption func(job *wrapJob)

func WithBeforeFunc(before BeforeFunc) JobOption {
	return func(job *wrapJob) {
		job.before = before
	}
}

func WithAfterFunc(after AfterFunc) JobOption {
	return func(job *wrapJob) {
		job.after = after
	}
}
