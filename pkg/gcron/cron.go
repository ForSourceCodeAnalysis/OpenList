package gcron

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var cronSchedule gocron.Scheduler

// Init 定时任务初始化
func Init() error {
	s, err := gocron.NewScheduler()
	if err != nil {
		logrus.Error(errors.WithStack(err))

		return err
	}
	cronSchedule = s
	cronSchedule.Start()
	return nil
}

// GetCronSchedule 获取定时任务
func GetCronSchedule() gocron.Scheduler {
	return cronSchedule
}
