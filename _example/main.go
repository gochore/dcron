package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gochore/dcron"
)

type RedisAtomic struct {
	client *redis.Client
}

func (m *RedisAtomic) SetIfNotExists(key, value string) bool {
	ret := m.client.SetNX(key, value, time.Hour)
	return ret.Err() == nil && ret.Val()
}

func main() {
	ra := &RedisAtomic{
		// init redis client
	}
	cron := dcron.NewCron(dcron.WithKey("TestCron"), dcron.WithAtomic(ra))

	job1 := dcron.NewJob("Job1", "*/15 * * * * *", func(ctx context.Context) error {
		if task, ok := dcron.TaskFromContext(ctx); ok {
			log.Println("run:", task.Job.Spec(), task.Key)
		}
		// do something
		return nil
	})
	job2 := dcron.NewJobWithAutoKey("*/20 * * * * *", Job2)
	job3 := Job3{}

	if err := cron.AddJobs(job1, job2, job3); err != nil {
		panic(err)
	}

	cron.Start()
	log.Println("cron started")
	time.Sleep(time.Minute)
	<-cron.Stop().Done()
	log.Println("cron stopped")
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
