package dcron

import (
	"reflect"
	"testing"
)

func TestStatistics_Add(t *testing.T) {
	type fields struct {
		TotalTask   int64
		PassedTask  int64
		FailedTask  int64
		SkippedTask int64
		MissedTask  int64
		TotalRun    int64
		PassedRun   int64
		FailedRun   int64
		RetriedRun  int64
	}
	type args struct {
		delta Statistics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Statistics
	}{
		{
			name: "regular",
			fields: fields{
				TotalTask:   1,
				PassedTask:  2,
				FailedTask:  3,
				SkippedTask: 4,
				MissedTask:  5,
				TotalRun:    6,
				PassedRun:   7,
				FailedRun:   8,
				RetriedRun:  9,
			},
			args: args{
				delta: Statistics{
					TotalTask:   1,
					PassedTask:  2,
					FailedTask:  3,
					SkippedTask: 4,
					MissedTask:  5,
					TotalRun:    6,
					PassedRun:   7,
					FailedRun:   8,
					RetriedRun:  9,
				},
			},
			want: Statistics{
				TotalTask:   2,
				PassedTask:  4,
				FailedTask:  6,
				SkippedTask: 8,
				MissedTask:  10,
				TotalRun:    12,
				PassedRun:   14,
				FailedRun:   16,
				RetriedRun:  18,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Statistics{
				TotalTask:   tt.fields.TotalTask,
				PassedTask:  tt.fields.PassedTask,
				FailedTask:  tt.fields.FailedTask,
				SkippedTask: tt.fields.SkippedTask,
				MissedTask:  tt.fields.MissedTask,
				TotalRun:    tt.fields.TotalRun,
				PassedRun:   tt.fields.PassedRun,
				FailedRun:   tt.fields.FailedRun,
				RetriedRun:  tt.fields.RetriedRun,
			}
			if got := s.Add(tt.args.delta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
