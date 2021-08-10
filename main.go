package main

import "blog/model"
import "blog/router"

func main() {
	model.Del()
	//model.Init()
	r := router.New()
	r.Run(":3000")
}
