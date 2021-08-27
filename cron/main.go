package cron

import (
	//"fmt"
	"github.com/robfig/cron/v3"
	"log"
	//"time"
)

func CronInit() {
	crontab := cron.New(cron.WithSeconds()) //精确到秒
	_, err := crontab.AddFunc(StatCronServiceSpec, StatCronService)
	if err != nil {
		panic(err)
	}
	crontab.Start()
	log.Println("定时任务以启动")
	select {}
}
