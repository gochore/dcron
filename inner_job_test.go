package dcron

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gochore/dcron/mock_dcron"
	"github.com/golang/mock/gomock"

	"github.com/robfig/cron/v3"
)

func Test_innerJob_Cron(t *testing.T) {
	c := NewCron()

	type fields struct {
		cron          *Cron
		entryID       cron.EntryID
		key           string
		spec          string
		before        BeforeFunc
		run           RunFunc
		after         AfterFunc
		retryTimes    int
		retryInterval RetryInterval
	}
	tests := []struct {
		name   string
		fields fields
		want   *Cron
	}{
		{
			name: "regular",
			fields: fields{
				cron: c,
			},
			want: c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &innerJob{
				cron:          tt.fields.cron,
				entryID:       tt.fields.entryID,
				key:           tt.fields.key,
				spec:          tt.fields.spec,
				before:        tt.fields.before,
				run:           tt.fields.run,
				after:         tt.fields.after,
				retryTimes:    tt.fields.retryTimes,
				retryInterval: tt.fields.retryInterval,
			}
			if got := j.Cron(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cron() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_innerJob_Key(t *testing.T) {
	type fields struct {
		cron          *Cron
		entryID       cron.EntryID
		key           string
		spec          string
		before        BeforeFunc
		run           RunFunc
		after         AfterFunc
		retryTimes    int
		retryInterval RetryInterval
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
			j := &innerJob{
				cron:          tt.fields.cron,
				entryID:       tt.fields.entryID,
				key:           tt.fields.key,
				spec:          tt.fields.spec,
				before:        tt.fields.before,
				run:           tt.fields.run,
				after:         tt.fields.after,
				retryTimes:    tt.fields.retryTimes,
				retryInterval: tt.fields.retryInterval,
			}
			if got := j.Key(); got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_innerJob_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEntryGetter := mock_dcron.NewMockentryGetter(ctrl)
	mutex := mock_dcron.NewMockMutex(ctrl)

	mockEntryGetter.EXPECT().
		Entry(gomock.Any()).
		DoAndReturn(func(id cron.EntryID) cron.Entry {
			now := time.Now()
			return cron.Entry{
				ID:   id,
				Next: now.Add(time.Duration(id) * time.Second),
				Prev: now,
			}
		}).
		MinTimes(1)

	mutex.EXPECT().
		SetIfNotExists(gomock.Any(), gomock.Any()).
		DoAndReturn(func(key, value string) bool {
			return value != "always_miss"
		}).
		MinTimes(1)

	type fields struct {
		cron          *Cron
		entryID       cron.EntryID
		entryGetter   entryGetter
		key           string
		spec          string
		before        BeforeFunc
		run           RunFunc
		after         AfterFunc
		retryTimes    int
		retryInterval RetryInterval
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "regular",
			fields: fields{
				cron:        NewCron(WithMutex(mutex)),
				entryID:     1,
				entryGetter: mockEntryGetter,
				before: func(task Task) (skip bool) {
					return false
				},
				run: func(ctx context.Context) error {
					return nil
				},
				after: func(task Task) {
					if task.Return != nil {
						t.Fatal(task.Return)
					}
				},
				retryTimes: 1,
			},
		},
		{
			name: "skip",
			fields: fields{
				cron:        NewCron(WithMutex(mutex)),
				entryID:     1,
				entryGetter: mockEntryGetter,
				before: func(task Task) (skip bool) {
					return true
				},
				run: func(ctx context.Context) error {
					return nil
				},
				after: func(task Task) {
					if !task.Skipped {
						t.Fatal(task.Skipped)
					}
				},
				retryTimes: 1,
			},
		},
		{
			name: "retry",
			fields: fields{
				cron:        NewCron(WithMutex(mutex)),
				entryID:     5,
				entryGetter: mockEntryGetter,
				before: func(task Task) (skip bool) {
					return false
				},
				run: func(ctx context.Context) error {
					return errors.New("show retry")
				},
				after: func(task Task) {
					if task.Return == nil {
						t.Fatal(task.Return)
					}
					if task.TriedTimes != 10 {
						t.Fatal(task.TriedTimes)
					}
				},
				retryTimes: 10,
			},
		},
		{
			name: "retry with interval",
			fields: fields{
				cron:        NewCron(WithMutex(mutex)),
				entryID:     5,
				entryGetter: mockEntryGetter,
				before: func(task Task) (skip bool) {
					return false
				},
				run: func(ctx context.Context) error {
					return errors.New("should retry")
				},
				after: func(task Task) {
					if task.Return == nil {
						t.Fatal(task.Return)
					}
					if task.TriedTimes >= 10 {
						t.Fatal(task.TriedTimes)
					}
				},
				retryTimes: 10,
				retryInterval: func(triedTimes int) time.Duration {
					return time.Duration(triedTimes) * time.Second
				},
			},
		},
		{
			name: "take too long",
			fields: fields{
				cron:        NewCron(WithMutex(mutex)),
				entryID:     1,
				entryGetter: mockEntryGetter,
				before: func(task Task) (skip bool) {
					return false
				},
				run: func(ctx context.Context) error {
					time.Sleep(2 * time.Second)
					return errors.New("show retry")
				},
				after: func(task Task) {
					if task.TriedTimes != 1 {
						t.Fatal(task.TriedTimes)
					}
					if task.Return == nil {
						t.Fatal(task.Return)
					}
				},
				retryTimes:    5,
				retryInterval: nil,
			},
		},
		{
			name: "miss",
			fields: fields{
				cron:        NewCron(WithMutex(mutex), WithHostname("always_miss")),
				entryID:     1,
				entryGetter: mockEntryGetter,
				before: func(task Task) (skip bool) {
					return false
				},
				run: func(ctx context.Context) error {
					return nil
				},
				after: func(task Task) {
					if !task.Missed {
						t.Fatal(task.Missed)
					}
				},
				retryTimes: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &innerJob{
				cron:          tt.fields.cron,
				entryID:       tt.fields.entryID,
				entryGetter:   tt.fields.entryGetter,
				key:           tt.fields.key,
				spec:          tt.fields.spec,
				before:        tt.fields.before,
				run:           tt.fields.run,
				after:         tt.fields.after,
				retryTimes:    tt.fields.retryTimes,
				retryInterval: tt.fields.retryInterval,
			}
			j.Run()
		})
	}
}

func Test_innerJob_Spec(t *testing.T) {
	type fields struct {
		cron          *Cron
		entryID       cron.EntryID
		key           string
		spec          string
		before        BeforeFunc
		run           RunFunc
		after         AfterFunc
		retryTimes    int
		retryInterval RetryInterval
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
			j := &innerJob{
				cron:          tt.fields.cron,
				entryID:       tt.fields.entryID,
				key:           tt.fields.key,
				spec:          tt.fields.spec,
				before:        tt.fields.before,
				run:           tt.fields.run,
				after:         tt.fields.after,
				retryTimes:    tt.fields.retryTimes,
				retryInterval: tt.fields.retryInterval,
			}
			if got := j.Spec(); got != tt.want {
				t.Errorf("Spec() = %v, want %v", got, tt.want)
			}
		})
	}
}
