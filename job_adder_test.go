package dcron

import (
	"testing"
)

func TestNewJobAdder(t *testing.T) {
	type args struct {
		cron *Cron
	}
	tests := []struct {
		name    string
		args    args
		jobs    []Job
		wantErr bool
	}{
		{
			name:    "regular",
			args:    args{
				cron: NewCron(),
			},
			jobs:    []Job{
				NewJob("test_job", "* * * * * *", nil),
			},
			wantErr: false,
		},
		{
			name:    "nil cron",
			args:    args{
				cron: nil,
			},
			jobs:    []Job{
				NewJob("test_job", "* * * * * *", nil),
			},
			wantErr: true,
		},
		{
			name:    "wrong spec",
			args:    args{
				cron: NewCron(),
			},
			jobs:    []Job{
				NewJob("test_job", "* * * * * +", nil),
			},
			wantErr: true,
		},
		{
			name:    "multiple wrong spec",
			args:    args{
				cron: NewCron(),
			},
			jobs:    []Job{
				NewJob("test_job2", "* * * * * +", nil),
				NewJob("test_job2", "* * * * *", nil),

			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adder := NewJobAdder(tt.args.cron)
			for _, job := range tt.jobs {
				adder.Add(job)
			}
			if err := adder.Err(); (err != nil) != tt.wantErr {
				t.Errorf("Err() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
