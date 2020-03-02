package dcron

import "context"

type JobMeta interface {
	Key() string
	Spec() string
}

type Job interface {
	JobMeta
	Run(ctx context.Context) error
	Options() []JobOption
}

type wrappedJob struct {
	key     string
	spec    string
	run     RunFunc
	options []JobOption
}

func NewJob(key, spec string, run RunFunc, options ...JobOption) Job {
	return &wrappedJob{
		key:     key,
		spec:    spec,
		run:     run,
		options: options,
	}
}

func (j *wrappedJob) Key() string {
	return j.key
}

func (j *wrappedJob) Spec() string {
	return j.spec
}

func (j *wrappedJob) Run(ctx context.Context) error {
	if j.run != nil {
		return j.run(ctx)
	}
	return nil
}

func (j *wrappedJob) Options() []JobOption {
	return j.options
}
