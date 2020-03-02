package dcron

import (
	"context"
	"time"
)

type JobOption func(job *innerJob)

type BeforeFunc func(task Task) (skip bool)

type RunFunc func(ctx context.Context) error

type AfterFunc func(task Task)

type RetryInterval func(triedTimes int) time.Duration

func WithBeforeFunc(before BeforeFunc) JobOption {
	return func(job *innerJob) {
		job.before = before
	}
}

func WithAfterFunc(after AfterFunc) JobOption {
	return func(job *innerJob) {
		job.after = after
	}
}

func WithRetryTimes(retryTimes int) JobOption {
	return func(job *innerJob) {
		job.retryTimes = retryTimes
	}
}

func WithRetryInterval(retryInterval RetryInterval) JobOption {
	return func(job *innerJob) {
		job.retryInterval = retryInterval
	}
}

func WithNoMutex() JobOption {
	return func(job *innerJob) {
		job.noMutex = true
	}
}
