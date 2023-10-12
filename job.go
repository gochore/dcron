package dcron

import (
	"context"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// Job describes a type which could be added to a cron.
type Job interface {
	// Key returns the unique key of the job.
	Key() string
	// Spec returns spec of the job, like "* * * * * *".
	Spec() string
	// Run is what the job do.
	Run(ctx context.Context) error
	// Options returns options of the job.
	Options() []JobOption
}

type wrappedJob struct {
	key     string
	spec    string
	run     RunFunc
	options []JobOption
}

// NewJob returns a new Job with specified options.
func NewJob(key, spec string, run RunFunc, options ...JobOption) Job {
	return &wrappedJob{
		key:     key,
		spec:    spec,
		run:     run,
		options: options,
	}
}

// NewJobWithAutoKey returns a new Job with the "run" function's name as key.
// Be careful, the "run" should be a non-anonymous function,
// or returned Job will have an emtpy key, and can not be added to a Cron.
func NewJobWithAutoKey(spec string, run RunFunc, options ...JobOption) Job {
	return NewJob(funcName(run), spec, run, options...)
}

// Key implements Job.Key.
func (j *wrappedJob) Key() string {
	return j.key
}

// Spec implements Job.Spec.
func (j *wrappedJob) Spec() string {
	return j.spec
}

// Run implements Job.Run.
func (j *wrappedJob) Run(ctx context.Context) error {
	if j.run != nil {
		return j.run(ctx)
	}
	return nil
}

// Options implements Job.Options.
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
