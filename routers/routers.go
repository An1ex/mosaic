package routers

import (
	"github.com/gin-gonic/gin"
	"mosaic/controller"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("./templates/*")

	r.GET("/", controller.Index)

	r.POST("/mosaic", controller.Mosaic)

	return r
}
