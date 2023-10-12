package dcron

import (
	"context"
	"reflect"
	"testing"
)

func TestNewJob(t *testing.T) {
	type args struct {
		key     string
		spec    string
		run     RunFunc
		options []JobOption
	}
	tests := []struct {
		name string
		args args
		want Job
	}{
		{
			name: "regular",
			args: args{
				key:     "test_job",
				spec:    "* * * * * *",
				run:     nil,
				options: nil,
			},
			want: &wrappedJob{
				key:     "test_job",
				spec:    "* * * * * *",
				run:     nil,
				options: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJob(tt.args.key, tt.args.spec, tt.args.run, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrappedJob_Key(t *testing.T) {
	type fields struct {
		key     string
		spec    string
		run     RunFunc
		options []JobOption
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "regular",
			fields: fields{
				key: "test_job",
			},
			want: "test_job",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &wrappedJob{
				key:     tt.fields.key,
				spec:    tt.fields.spec,
				run:     tt.fields.run,
				options: tt.fields.options,
			}
			if got := j.Key(); got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrappedJob_Options(t *testing.T) {
	options := []JobOption{WithRetryTimes(1)}

	type fields struct {
		key     string
		spec    string
		run     RunFunc
		options []JobOption
	}
	tests := []struct {
		name   string
		fields fields
		want   []JobOption
	}{
		{
			name: "regular",
			fields: fields{
				options: options,
			},
			want: options,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &wrappedJob{
				key:     tt.fields.key,
				spec:    tt.fields.spec,
				run:     tt.fields.run,
				options: tt.fields.options,
			}
			if got := j.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrappedJob_Run(t *testing.T) {
	type fields struct {
		key     string
		spec    string
		run     RunFunc
		options []JobOption
	}
	type args struct {
		ctx context.Context
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
				run: func(ctx context.Context) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "nil run",
			fields: fields{
				run: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &wrappedJob{
				key:     tt.fields.key,
				spec:    tt.fields.spec,
				run:     tt.fields.run,
				options: tt.fields.options,
			}
			if err := j.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_wrappedJob_Spec(t *testing.T) {
	type fields struct {
		key     string
		spec    string
		run     RunFunc
		options []JobOption
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "regular",
			fields: fields{
				spec: "* * * * * *",
			},
			want: "* * * * * *",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &wrappedJob{
				key:     tt.fields.key,
				spec:    tt.fields.spec,
				run:     tt.fields.run,
				options: tt.fields.options,
			}
			if got := j.Spec(); got != tt.want {
				t.Errorf("Spec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJobWithAutoKey(t *testing.T) {
	type args struct {
		spec    string
		run     RunFunc
		options []JobOption
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, job Job)
	}{
		{
			name: "regular",
			args: args{
				run: Func,
			},
			check: func(t *testing.T, job Job) {
				if job.Key() != "Func" {
					t.Fatal(job.Key())
				}
			},
		},
		{
			name: "method",
			args: args{
				run: Foo{}.Method,
			},
			check: func(t *testing.T, job Job) {
				if job.Key() != "Method" {
					t.Fatal(job.Key())
				}
			},
		},
		{
			name: "nil",
			args: args{
				run: nil,
			},
			check: func(t *testing.T, job Job) {
				if job.Key() != "" {
					t.Fatal(job.Key())
				}
			},
		},
		{
			name: "anonymous func",
			args: args{
				run: func(ctx context.Context) error {
					return nil
				},
			},
			check: func(t *testing.T, job Job) {
				if job.Key() != "" {
					t.Fatal(job.Key())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewJobWithAutoKey(tt.args.spec, tt.args.run, tt.args.options...)
			tt.check(t, got)
		})
	}
}

func Func(_ context.Context) error {
	return nil
}

type Foo struct {
}

func (Foo) Method(_ context.Context) error {
	return nil
}
