package main

// @title 博客系统后端
// @version 1.0
// @description 完成了博客的基本功能

// @contact.name 陈浩
// @contact.url https://github.com/chen102
// @contact.email 773532732@qq.com

// @host localhost
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
