package dcron

import (
	"testing"

	"github.com/gochore/dcron/mock_dcron"
)

func TestWithHostname(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option CronOption)
	}{
		{
			name: "regular",
			args: args{
				hostname: "test_hostname",
			},
			check: func(t *testing.T, option CronOption) {
				c := NewCron()
				option(c)
				if c.hostname != "test_hostname" {
					t.Fatal(c.hostname)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithHostname(tt.args.hostname)
			tt.check(t, got)
		})
	}
}

func TestWithKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option CronOption)
	}{
		{
			name: "regular",
			args: args{
				key: "test_cron",
			},
			check: func(t *testing.T, option CronOption) {
				c := NewCron()
				option(c)
				if c.key != "test_cron" {
					t.Fatal(c.key)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithKey(tt.args.key)
			tt.check(t, got)
		})
	}
}

func TestWithAtomic(t *testing.T) {
	type args struct {
		atomic Atomic
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option CronOption)
	}{
		{
			name: "regular",
			args: args{
				atomic: mock_dcron.NewMockAtomic(nil),
			},
			check: func(t *testing.T, option CronOption) {
				c := NewCron()
				option(c)
				if c.atomic == nil {
					t.Fatal(c.atomic)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithAtomic(tt.args.atomic)
			tt.check(t, got)
		})
	}
}
