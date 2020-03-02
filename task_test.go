package dcron

import (
	"context"
	"reflect"
	"testing"
)

func TestTaskInContext(t *testing.T) {
	task := Task{
		Key: "test_task",
	}
	zero := Task{}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  Task
		want1 bool
	}{
		{
			name: "regular",
			args: args{
				ctx: context.WithValue(context.Background(), keyContextTask, task),
			},
			want:  task,
			want1: true,
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
			},
			want:  zero,
			want1: false,
		},
		{
			name: "context without task",
			args: args{
				ctx: context.Background(),
			},
			want:  zero,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := TaskFromContext(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskFromContext() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TaskFromContext() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
