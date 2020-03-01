package dcron

import (
	"log"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func Test_Cron(t *testing.T) {
	cron := NewCron()
	if err := cron.AddJob(&TestJob{
		key:  "test",
		spec: "*/5 * * * * *",
	}); err != nil {
		t.Fatal(err)
	}
	cron.Start()
	time.Sleep(10 * time.Second)
	<-cron.Stop().Done()
}

type TestJob struct {
	key  string
	spec string
}

func (j *TestJob) Key() string {
	return j.key
}

func (j *TestJob) Spec() string {
	return j.spec
}

func (j *TestJob) Before(ctx Context) (skip bool) {
	log.Println("before")
	spew.Dump(ctx)
	return false
}

func (j *TestJob) Run() error {
	log.Println("run")
	return nil
}

func (j *TestJob) After(ctx Context) {
	log.Println("after")
	spew.Dump(ctx)
}
