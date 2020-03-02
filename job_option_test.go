package dcron

import (
	"fmt"
	"testing"
	"time"
)

func TestWithAfterFunc(t *testing.T) {
	after := func(task Task) {

	}

	type args struct {
		after AfterFunc
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option JobOption)
	}{
		{
			name: "regular",
			args: args{
				after: after,
			},
			check: func(t *testing.T, option JobOption) {
				j := &innerJob{}
				option(j)
				if fmt.Sprintf("%p", j.after) != fmt.Sprintf("%p", after) {
					t.Fatal()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithAfterFunc(tt.args.after)
			tt.check(t, got)
		})
	}
}

func TestWithBeforeFunc(t *testing.T) {
	before := func(task Task) (skip bool) {
		return false
	}

	type args struct {
		before BeforeFunc
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option JobOption)
	}{
		{
			name: "regular",
			args: args{
				before: before,
			},
			check: func(t *testing.T, option JobOption) {
				j := &innerJob{}
				option(j)
				if fmt.Sprintf("%p", j.before) != fmt.Sprintf("%p", before) {
					t.Fatal()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithBeforeFunc(tt.args.before)
			tt.check(t, got)
		})
	}
}

func TestWithRetryInterval(t *testing.T) {
	retryInterval := func(triedTimes int) time.Duration {
		return time.Second
	}
	type args struct {
		retryInterval RetryInterval
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option JobOption)
	}{
		{
			name: "regular",
			args: args{
				retryInterval: retryInterval,
			},
			check: func(t *testing.T, option JobOption) {
				j := &innerJob{}
				option(j)
				if fmt.Sprintf("%p", j.retryInterval) != fmt.Sprintf("%p", retryInterval) {
					t.Fatal()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithRetryInterval(tt.args.retryInterval)
			tt.check(t, got)
		})
	}
}

func TestWithRetryTimes(t *testing.T) {
	type args struct {
		retryTimes int
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option JobOption)
	}{
		{
			name: "regular",
			args: args{
				retryTimes: 10,
			},
			check: func(t *testing.T, option JobOption) {
				j := &innerJob{}
				option(j)
				if j.retryTimes != 10 {
					t.Fatal(j.retryTimes)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithRetryTimes(tt.args.retryTimes)
			tt.check(t, got)
		})
	}
}
