package routers

import (
	"github.com/gin-gonic/gin"
	themsController "../controller"
	"database/sql"
)

func Init(db *sql.DB) *gin.Engine {
	r := gin.New()

	r.LoadHTMLGlob("www/*.html")
	r.Static("/assets", "./www/assets")

	baseController := themsController.NewBaseController(db)

	r.GET("/", themsController.IndexController{baseController}.GetIndex)

	// ログイン関連
	r.GET("/login", themsController.LoginController{baseController}.GetLogin)
	r.POST("/login", themsController.LoginController{baseController}.PostLogin)

	return r
}
