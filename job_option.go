package dcron

import (
	"context"
	"time"
)

// JobOption represents a modification to the default behavior of a Job.
type JobOption func(job *innerJob)

// BeforeFunc represents the function could be called before Run.
type BeforeFunc func(task Task) (skip bool)

// RunFunc represents the function could be called by a cron.
type RunFunc func(ctx context.Context) error

// AfterFunc represents the function could be called after Run.
type AfterFunc func(task Task)

// RetryInterval indicates how long should delay before retrying when run failed `triedTimes` times.
type RetryInterval func(triedTimes int) time.Duration

// WithBeforeFunc specifies what to do before Run.
func WithBeforeFunc(before BeforeFunc) JobOption {
	return func(job *innerJob) {
		job.before = before
	}
}

// WithAfterFunc specifies what to do after Run.
func WithAfterFunc(after AfterFunc) JobOption {
	return func(job *innerJob) {
		job.after = after
	}
}

// WithRetryTimes specifies max times to retry,
// retryTimes will be set as 1 if it is less than 1.
func WithRetryTimes(retryTimes int) JobOption {
	return func(job *innerJob) {
		job.retryTimes = retryTimes
	}
}

// WithRetryInterval indicates how long should delay before retrying when run failed `triedTimes` times.
func WithRetryInterval(retryInterval RetryInterval) JobOption {
	return func(job *innerJob) {
		job.retryInterval = retryInterval
	}
}

// WithNoMutex means the job will run at multiple cron instances,
// even though the cron has Atomic.
func WithNoMutex() JobOption {
	return func(job *innerJob) {
		job.noMutex = true
	}
}

// WithGroup adds the current job to the group.
func WithGroup(group Group) JobOption {
	return func(job *innerJob) {
		job.group = group
	}
}
