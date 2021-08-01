package main

import "blog/model"

func main() {
	model.Del()
	model.Init()
	r := New()
	r.Run(":3000")
}
