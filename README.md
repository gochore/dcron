# dcron

[![Go Reference](https://pkg.go.dev/badge/github.com/gochore/dcron.svg)](https://pkg.go.dev/github.com/gochore/dcron)
[![Actions](https://github.com/gochore/dcron/actions/workflows/test.yaml/badge.svg)](https://github.com/gochore/dcron/actions)
[![Codecov](https://codecov.io/gh/gochore/dcron/branch/master/graph/badge.svg)](https://codecov.io/gh/gochore/dcron)
[![Go Report Card](https://goreportcard.com/badge/github.com/gochore/dcron)](https://goreportcard.com/report/github.com/gochore/dcron)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gochore/dcron)](https://github.com/gochore/dcron/blob/master/go.mod)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/gochore/dcron)](https://github.com/gochore/dcron/releases)

A distributed cron framework.

## Install

```shell
go get github.com/gochore/dcron
```

## Quick Start

First, implement a distributed atomic operation that only requires support for one method: `SetIfNotExists`.
You can implement it in any way you prefer, such as using Redis `SetNX`.

```go
import "github.com/redis/go-redis/v9"

type RedisAtomic struct {
	client *redis.Client
}

func (m *RedisAtomic) SetIfNotExists(ctx context.Context, key, value string) bool {
	ret := m.client.SetNX(ctx, key, value, time.Hour)
	return ret.Err() == nil && ret.Val()
}
```

Now you can create a cron with that:

```go
func main() {
	atomic := &RedisAtomic{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
	cron := dcron.NewCron(dcron.WithKey("TestCron"), dcron.WithAtomic(atomic))
}
```

Then, create a job and add it to the cron.

```go
	job1 := dcron.NewJob("Job1", "*/15 * * * * *", func(ctx context.Context) error {
		if task, ok := dcron.TaskFromContext(ctx); ok {
			log.Println("run:", task.Job.Spec(), task.Key)
		}
		// do something
		return nil
	})
	if err := cron.AddJobs(job1); err != nil {
		log.Fatal(err)
	}
```

Finally, start the cron:

```go
	cron.Start()
	log.Println("cron started")
	time.Sleep(time.Minute)
	<-cron.Stop().Done()
```

If you start the program multiple times, you will notice that the cron will run the job once every 15 seconds on only one of the processes.

| process 1                                                              | process 2                                                              | process 3                                                              |
|------------------------------------------------------------------------|------------------------------------------------------------------------|------------------------------------------------------------------------|
| 2023/10/13 11:39:45 cron started                                       | 2023/10/13 11:39:47 cron started                                       | 2023/10/13 11:39:48 cron started                                       |
|                                                                        |                                                                        | 2023/10/13 11:40:00 run: */15 * * * * * dcron:TestCron.Job1@1697168400 |
|                                                                        | 2023/10/13 11:40:15 run: */15 * * * * * dcron:TestCron.Job1@1697168415 |                                                                        |
|                                                                        |                                                                        | 2023/10/13 11:40:30 run: */15 * * * * * dcron:TestCron.Job1@1697168430 |
| 2023/10/13 11:40:45 run: */15 * * * * * dcron:TestCron.Job1@1697168445 |                                                                        |                                                                        |
