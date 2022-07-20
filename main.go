package main

import (
	"mosaic/routers"
)

func main() {
	r := routers.SetUpRouter()
	err := r.Run()
	if err != nil {
		return
	}
}
