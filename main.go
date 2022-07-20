package main

import (
	"mosaic/model/tiles"
	"mosaic/routers"
)

func modelInit() {
	tiles.Init()
}
func main() {
	modelInit()

	r := routers.SetUpRouter()
	err := r.Run()
	if err != nil {
		return
	}
}
