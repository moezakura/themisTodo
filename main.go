package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.New()

	r.LoadHTMLGlob("www/*.html")
	r.Static("/assets", "./www/assets")

	r.GET("/", GetRoot)
	r.Run(":31204")
}

func GetRoot(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", gin.H{})
}