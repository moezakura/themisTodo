package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	themsController "../controller"
	"database/sql"
)

func Init(db *sql.DB) *gin.Engine {
	r := gin.New()

	r.LoadHTMLGlob("www/*.html")
	r.Static("/fontawesome", "./www/assets/fontawesome/web-fonts-with-css/")
	r.Static("/assets", "./www/assets")
	r.HTMLRender = InitRender()

	baseController := themsController.NewBaseController(db, r)

	r.GET("/", themsController.IndexController{baseController}.GetIndex)

	// ログイン関連
	r.GET("/login", themsController.LoginController{baseController}.GetLogin)
	r.POST("/login", themsController.LoginController{baseController}.PostLogin)

	// メイン画面
	r.GET("/home", themsController.HomeController{baseController}.GetHome)

	// プロジェクト関連
	projects := r.Group("/project")
	{
		projects.GET("/add", themsController.ProjectsController{baseController}.GetAdd)
		projects.POST("/add", themsController.ProjectsController{baseController}.PostAdd)

		projects.GET("/view/:projectId", themsController.ProjectsController{baseController}.GetTaskBoard)
		projects.POST("/view/:projectId", themsController.ProjectsController{baseController}.PostTaskBoard)
	}

	return r
}

func InitRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("index", "www/base.html", "www/index.html")
	r.AddFromFiles("login", "www/base.html", "www/login.html")
	r.AddFromFiles("home", "www/base.html", "www/header.html", "www/home.html")
	r.AddFromFiles("projectAdd", "www/base.html", "www/header.html", "www/projectAdd.html")
	r.AddFromFiles("projectTaskBoard", "www/base.html", "www/header.html", "www/taskBoard.html")

	return r
}
