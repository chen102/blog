package main

import "blog/model"
import "blog/router"
import "blog/cron"
import "blog/tool"

func main() {
	config := tool.NewConfig()
	config.ReadConfig()
	model.DelMysql(*config)
	model.DelRedis(*config)
	go cron.CronInit()
	r := router.New()
	r.Run(":3000")
}
