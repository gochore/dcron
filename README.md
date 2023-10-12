# dcron

[![Go Reference](https://pkg.go.dev/badge/github.com/gochore/dcron.svg)](https://pkg.go.dev/github.com/gochore/dcron)
[![Build Status](https://travis-ci.com/gochore/dcron.svg?branch=master)](https://travis-ci.com/gochore/dcron)
[![codecov](https://codecov.io/gh/gochore/dcron/branch/master/graph/badge.svg)](https://codecov.io/gh/gochore/dcron)
[![Go Report Card](https://goreportcard.com/badge/github.com/gochore/dcron)](https://goreportcard.com/report/github.com/gochore/dcron)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gochore/dcron)](https://github.com/gochore/dcron/blob/master/go.mod)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/gochore/dcron)](https://github.com/gochore/dcron/releases)

A distributed cron framework.

## Install

```shell
go get github.com/gochore/dcron
```

## Example

First of all, you should implement a distributed atomic operation:

```go
type Atomic interface {
	SetIfNotExists(key, value string) bool
}
```

You can implement it any way you like, for example, via Redis `SetNX`:

```go
type RedisAtomic struct {
	client *redis.Client
}

func (m *RedisAtomic) SetIfNotExists(key, value string) bool {
	ret := m.client.SetNX(key, value, time.Hour)
	return ret.Err() == nil && ret.Val()
}
```

Now we can create a cron with that:

```go
	ra := &RedisAtomic{
		// init redis client
	}
	cron := dcron.NewCron(dcron.WithKey("TestCron"), dcron.WithAtomic(ra))
```

There are many ways to create jobs:
- use `dcron.NewJob`;
- use `dcron.NewJobWithAutoKey`;
- implement interface `dcron.Job`.

```go
	
func main {
	job1 := dcron.NewJob("Job1", "*/15 * * * * *", func(ctx context.Context) error {
		if task, ok := dcron.TaskFromContext(ctx); ok {
			log.Println("run:", task.Job.Spec(), task.Key)
		}
		// do something
		return nil
	})
	job2 := dcron.NewJobWithAutoKey("*/20 * * * * *", Job2)
	job3 := Job3{}
}

func Job2(ctx context.Context) error {
	if task, ok := dcron.TaskFromContext(ctx); ok {
		log.Println("run:", task.Job.Spec(), task.Key)
	}
	// do something
	return nil
}

type Job3 struct {
}

func (j Job3) Key() string {
	return "Job3"
}

func (j Job3) Spec() string {
	return "*/30 * * * * *"
}

func (j Job3) Run(ctx context.Context) error {
	if task, ok := dcron.TaskFromContext(ctx); ok {
		log.Println("run:", task.Job.Spec(), task.Key)
	}
	// do something
	return nil
}

func (j Job3) Options() []dcron.JobOption {
	return nil
}
```

Finally, add the jobs to the cron, and start it:

```go
	if err := cron.AddJobs(job1, job2, job3); err != nil {
		panic(err)
	}

	cron.Start()
	log.Println("cron started")
	time.Sleep(time.Minute)
	<-cron.Stop().Done()
	log.Println("cron stopped")
```

You will see logging:

```text
2020/03/11 15:28:04 cron started
2020/03/11 15:28:15 run: */15 * * * * * dcron:TestCron.Job1@1583911695
2020/03/11 15:28:20 run: */20 * * * * * dcron:TestCron.Job2@1583911700
2020/03/11 15:28:30 run: */30 * * * * * dcron:TestCron.Job3@1583911710
2020/03/11 15:28:30 run: */15 * * * * * dcron:TestCron.Job1@1583911710
2020/03/11 15:28:40 run: */20 * * * * * dcron:TestCron.Job2@1583911720
2020/03/11 15:28:45 run: */15 * * * * * dcron:TestCron.Job1@1583911725
2020/03/11 15:29:00 run: */30 * * * * * dcron:TestCron.Job3@1583911740
2020/03/11 15:29:00 run: */15 * * * * * dcron:TestCron.Job1@1583911740
2020/03/11 15:29:00 run: */20 * * * * * dcron:TestCron.Job2@1583911740
2020/03/11 15:29:04 cron stopped
```

There is the complete example: [_example/main.go](https://github.com/gochore/dcron/tree/_example/main.go).
