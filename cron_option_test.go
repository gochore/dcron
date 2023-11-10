package dcron

import (
	"context"
	"testing"
	"time"

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

func TestWithLocation(t *testing.T) {
	type args struct {
		loc *time.Location
	}
	tests := []struct {
		name  string
		args  args
		want  CronOption
		check func(t *testing.T, option CronOption)
	}{
		{
			name: "regular",
			args: args{
				loc: time.UTC,
			},
			check: func(t *testing.T, option CronOption) {
				c := NewCron()
				option(c)
				if c.location != time.UTC {
					t.Fatal(c.location)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithLocation(tt.args.loc)
			tt.check(t, got)
		})
	}
}

func TestWithContext(t *testing.T) {
	type ctxKey struct{}
	job := NewJob("test", "*/2 * * * * *", func(ctx context.Context) error {
		task, _ := TaskFromContext(ctx)
		t.Logf("run: %v %v", task.Job.Key(), task.PlanAt.Format(time.RFC3339))
		if ctx.Value(ctxKey{}).(string) != "test" {
			t.Fatal("wrong context")
		}
		return nil
	})

	t.Run("start with context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, ctxKey{}, "test")

		c := NewCron(WithContext(ctx))
		if err := c.AddJobs(job); err != nil {
			t.Fatal(err)
		}
		c.Start()
		time.Sleep(5 * time.Second)
		<-c.Stop().Done()
	})

	t.Run("run with context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, ctxKey{}, "test")

		c := NewCron(WithContext(ctx))
		if err := c.AddJobs(job); err != nil {
			t.Fatal(err)
		}
		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()
		c.Run()
	})
}
