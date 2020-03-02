package dcron

import (
	"context"
	"testing"
	"time"

	"github.com/gochore/dcron/mock_dcron"
	"github.com/golang/mock/gomock"
	"github.com/robfig/cron/v3"
)

func Test_Cron(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mutex := mock_dcron.NewMockMutex(ctrl)

	c := NewCron(WithKey("test_cron"), WithMutex(mutex))

	mutex.EXPECT().
		SetIfNotExists(gomock.Any(), c.Hostname()).
		Return(true).
		Times(2)

	job := NewJob("test", "*/5 * * * * *", func(ctx context.Context) error {
		task, _ := TaskFromContext(ctx)
		select {
		case <-ctx.Done():
			t.Logf("exit: %+v", task)
		case <-time.After(time.Second):
			t.Logf("run: %+v", task)
		}
		return nil
	}, WithBeforeFunc(func(task Task) (skip bool) {
		t.Logf("before: %+v", task)
		return false
	}), WithAfterFunc(func(task Task) {
		t.Logf("after: %+v", task)
	}))
	if err := c.AddJobs(job); err != nil {
		t.Fatal(err)
	}
	c.Start()
	c.Run() // should be not working
	time.Sleep(10 * time.Second)
	<-c.Stop().Done()
}

func TestCron_AddJobs(t *testing.T) {
	c := cron.New(cron.WithSeconds())

	type fields struct {
		key      string
		hostname string
		cron     *cron.Cron
		mutex    Mutex
		jobs     []*innerJob
	}
	type args struct {
		jobs []Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "regular",
			fields: fields{
				cron: c,
			},
			args: args{
				jobs: []Job{
					NewJob("test_job", "* * * * * *", nil),
				},
			},
			wantErr: false,
		},
		{
			name: "multiple jobs",
			fields: fields{
				cron: c,
			},
			args: args{
				jobs: []Job{
					NewJob("test_job_1", "* * * * * *", nil),
					NewJob("test_job_2", "* * * * * *", nil),
					NewJob("test_job_3", "* * * * * *", nil),
				},
			},
			wantErr: false,
		},
		{
			name: "multiple jobs contain error",
			fields: fields{
				cron: c,
			},
			args: args{
				jobs: []Job{
					NewJob("test_job_1", "* * * * * *", nil),
					NewJob("test_job_2", "* * * * *", nil),
					NewJob("test_job_3", "* * * * * *", nil),
				},
			},
			wantErr: true,
		},
		{
			name: "multiple jobs contain same",
			fields: fields{
				cron: c,
			},
			args: args{
				jobs: []Job{
					NewJob("test_job_1", "* * * * * *", nil),
					NewJob("test_job_2", "* * * * * *", nil),
					NewJob("test_job_1", "* * * * * *", nil),
				},
			},
			wantErr: true,
		},
		{
			name: "with option",
			fields: fields{
				cron: c,
			},
			args: args{
				jobs: []Job{
					NewJob("test_job", "* * * * * *", nil, WithRetryTimes(3)),
				},
			},
			wantErr: false,
		},
		{
			name: "wrong spec",
			fields: fields{
				cron: c,
			},
			args: args{
				jobs: []Job{
					NewJob("test_job", "* * * * *", nil),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cron{
				key:      tt.fields.key,
				hostname: tt.fields.hostname,
				cron:     tt.fields.cron,
				mutex:    tt.fields.mutex,
				jobs:     tt.fields.jobs,
			}
			if err := c.AddJobs(tt.args.jobs...); (err != nil) != tt.wantErr {
				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCron_Hostname(t *testing.T) {
	type fields struct {
		key      string
		hostname string
		cron     *cron.Cron
		mutex    Mutex
		jobs     []*innerJob
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "regular",
			fields: fields{
				hostname: "test_hostname",
			},
			want: "test_hostname",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cron{
				key:      tt.fields.key,
				hostname: tt.fields.hostname,
				cron:     tt.fields.cron,
				mutex:    tt.fields.mutex,
				jobs:     tt.fields.jobs,
			}
			if got := c.Hostname(); got != tt.want {
				t.Errorf("Hostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCron_Key(t *testing.T) {
	type fields struct {
		key      string
		hostname string
		cron     *cron.Cron
		mutex    Mutex
		jobs     []*innerJob
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "regular",
			fields: fields{
				key: "test_key",
			},
			want: "test_key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cron{
				key:      tt.fields.key,
				hostname: tt.fields.hostname,
				cron:     tt.fields.cron,
				mutex:    tt.fields.mutex,
				jobs:     tt.fields.jobs,
			}
			if got := c.Key(); got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCron(t *testing.T) {
	type args struct {
		options []CronOption
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, c *Cron)
	}{
		{
			name: "regular",
			args: args{
				options: nil,
			},
			check: func(t *testing.T, c *Cron) {
				if c == nil {
					t.Fatal(t)
				}
			},
		},
		{
			name: "with_option",
			args: args{
				options: []CronOption{WithKey("test_cron")},
			},
			check: func(t *testing.T, c *Cron) {
				if c == nil {
					t.Fatal(t)
				}
				if c.key != "test_cron" {
					t.Fatal(c.key)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCron(tt.args.options...)
			tt.check(t, got)
		})
	}
}

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

func TestWithMutex(t *testing.T) {
	type args struct {
		mutex Mutex
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, option CronOption)
	}{
		{
			name: "regular",
			args: args{
				mutex: mock_dcron.NewMockMutex(nil),
			},
			check: func(t *testing.T, option CronOption) {
				c := NewCron()
				option(c)
				if c.mutex == nil {
					t.Fatal(c.mutex)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithMutex(tt.args.mutex)
			tt.check(t, got)
		})
	}
}
