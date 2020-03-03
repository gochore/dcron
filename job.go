package dcron

import (
	"context"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

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

// NewJobWithNonAnonymousFunc return a new Job with the 'run' function's name as key.
// Be careful, the 'run' should be a non-anonymous function,
// or returned Job will has a emtpy key, and can not be added to a Cron.
func NewJobWithNonAnonymousFunc(spec string, run RunFunc, options ...JobOption) Job {
	return &wrappedJob{
		key:     funcName(run),
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

func funcName(run RunFunc) string {
	if run != nil {
		name := runtime.FuncForPC(reflect.ValueOf(run).Pointer()).Name()
		splits := strings.Split(name, ".")
		name = strings.TrimSuffix(splits[len(splits)-1], "-fm") // method closures have a "-fm" suffix
		if regexp.MustCompile("^func[0-9]+$").MatchString(name) {
			return ""
		}
		return name
	}
	return ""
}
