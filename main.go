package main

import "blog/model"
import "blog/router"

func main() {
	if err := model.DelMysql(); err != nil {
		panic(err)
	}
	model.DB.AutoMigrate(&model.User{})
	model.DelRedis()
	r := router.New()
	r.Run(":3000")
}
