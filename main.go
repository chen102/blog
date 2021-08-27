package main

import "blog/model"
import "blog/router"
import "blog/cron"

func main() {
	if err := model.DelMysql(); err != nil {
		panic(err)
	}
	go cron.CronInit()

	model.DelRedis()
	r := router.New()
	r.Run(":3000")
}
