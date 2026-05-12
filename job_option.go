package dcron

import (
	"context"
	"time"
)

// JobOption represents a modification to the default behavior of a Job.
type JobOption func(job *innerJob)

// BeforeFunc represents the function could be called before Run.
// Deprecated: Use BeforeContextFunc instead, it will be ignored if BeforeContextFunc is set.
type BeforeFunc func(task Task) (skip bool)

// BeforeContextFunc represents the function could be called before Run with the given Task.
type BeforeContextFunc func(ctx context.Context, task Task) (skip bool)

// RunFunc represents the function could be called by a cron.
type RunFunc func(ctx context.Context) error

// AfterFunc represents the function could be called after Run.
// Deprecated: Use AfterContextFunc instead, it will be ignored if AfterContextFunc is set.
type AfterFunc func(task Task)

// AfterContextFunc represents the function could be called after Run with the given Task.
type AfterContextFunc func(ctx context.Context, task Task)

// RetryInterval indicates how long should delay before retrying when run failed `triedTimes` times.
type RetryInterval func(triedTimes int) time.Duration

// DeriveContext indicates how to derive a new context from the job's base context and the current Task.
type DeriveContext func(ctx context.Context, task Task) context.Context

// WithBeforeFunc specifies what to do before Run.
// Deprecated: Use WithBeforeContextFunc instead.
func WithBeforeFunc(before BeforeFunc) JobOption {
	return func(job *innerJob) {
		job.before = before
	}
}

// WithBeforeContextFunc specifies what to do before Run with the given Task.
func WithBeforeContextFunc(before BeforeContextFunc) JobOption {
	return func(job *innerJob) {
		job.ctxBefore = before
	}
}

// WithAfterFunc specifies what to do after Run.
// Deprecated: Use WithAfterContextFunc instead.
func WithAfterFunc(after AfterFunc) JobOption {
	return func(job *innerJob) {
		job.after = after
	}
}

// WithAfterContextFunc specifies what to do after Run with the given Task.
func WithAfterContextFunc(after AfterContextFunc) JobOption {
	return func(job *innerJob) {
		job.ctxAfter = after
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

// WithDeriveContext specifies how to derive a new context for the entire job execution, including
// before/after hooks, Run, and retry logic. The returned context must derive from the provided ctx
// to preserve the deadline, cancellation signal, and the embedded Task value (accessible via TaskFromContext).
// Returning a detached context (e.g., context.Background()) will break deadline enforcement and retry timeout logic.
func WithDeriveContext(deriveContext DeriveContext) JobOption {
	return func(job *innerJob) {
		job.deriveContext = deriveContext
	}
}
