package dcron

import (
	"log"
	"testing"
	"time"

	"github.com/gochore/dcron/mock_dcron"
	"github.com/golang/mock/gomock"
)

func Test_Cron(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mutex := mock_dcron.NewMockMutex(ctrl)
	mutex.EXPECT().
		SetIfNotExists(gomock.Any(), gomock.Any()).
		Return(true).
		MinTimes(1)

	cron := NewCron(WithKey("test_cron"), WithMutex(mutex))
	job := NewJob("test", "*/5 * * * * *", func(task Task) error {
		log.Println("run")
		log.Printf("%+v\n", task)
		return nil
	}, WithBeforeFunc(func(task Task) (skip bool) {
		log.Println("before")
		log.Printf("%+v\n", task)
		return false
	}), WithAfterFunc(func(task Task) {
		log.Println("after")
		log.Printf("%+v\n", task)
	}))
	if err := cron.AddJob(job); err != nil {
		t.Fatal(err)
	}
	cron.Start()
	time.Sleep(10 * time.Second)
	<-cron.Stop().Done()
}
