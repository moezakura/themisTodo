package routers

import (
	"github.com/gin-gonic/gin"
	themsController "../controller"
)

func Init() *gin.Engine {
	r := gin.New()

	r.LoadHTMLGlob("www/*.html")
	r.Static("/assets", "./www/assets")

	r.GET("/", themsController.IndexController{}.GetIndex)

	return r
}
