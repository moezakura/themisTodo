package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	themsController "../controller"
	"database/sql"
	"html/template"
)

func Init(db *sql.DB) *gin.Engine {
	r := gin.New()

	r.LoadHTMLGlob("www/*.html")
	r.Static("/fontawesome", "./www/assets/fontawesome/web-fonts-with-css/")
	r.Static("/assets", "./www/assets")
	r.SetFuncMap(InitRenderFunc())
	r.HTMLRender = InitRender()

	baseController := themsController.NewBaseController(db, r)

	r.GET("/", themsController.IndexController{baseController}.GetIndex)

	// ログイン関連
	r.GET("/login", themsController.LoginController{baseController}.GetLogin)
	r.POST("/login", themsController.LoginController{baseController}.PostLogin)

	// メイン画面
	r.GET("/home", themsController.HomeController{baseController}.GetHome)

	//マイページ設定
	r.GET("/settings", themsController.HomeController{baseController}.GetSettings)

	// プロジェクト関連
	projects := r.Group("/project")
	{
		projects.GET("/add", themsController.ProjectsController{baseController}.GetAdd)
		projects.POST("/add", themsController.ProjectsController{baseController}.PostAdd)
		projects.POST("/delete/:projectId", themsController.ProjectsController{baseController}.PostDeleteProject)

		projects.GET("/view/:projectId", themsController.ProjectsController{baseController}.GetTaskBoard)

		projects.POST("/update/:projectId", themsController.ProjectsController{baseController}.PostUpdate)
		projects.POST("/addUser/:projectId", themsController.ProjectsController{baseController}.PostAddUser)
	}

	// タスク管理
	tasks := r.Group("/tasks")
	{
		tasks.POST("/create", themsController.TasksController{baseController}.PostTaskCreate)
		tasks.POST("/update/:createDate", themsController.TasksController{baseController}.PostUpdate)
		tasks.GET("/view/:createDate", themsController.TasksController{baseController}.GetView)
	}

	//アカウント関連
	account := r.Group("/account")
	{
		account.GET("/add", themsController.AccountController{baseController}.GetAdd)
		account.POST("/add", themsController.AccountController{baseController}.PostAdd)
		account.GET("/search", themsController.AccountController{baseController}.GetSearch)
		account.POST("/update/:accountUuid", themsController.AccountController{baseController}.PostUpdate)
		account.POST("/updateIcon/:accountUuid", themsController.AccountController{baseController}.PostUpdateIcon)
	}
	return r
}

func InitRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("index", "www/base.html", "www/index.html")
	r.AddFromFiles("login", "www/base.html", "www/login.html")
	r.AddFromFiles("home", "www/base.html", "www/header.html", "www/home.html")
	r.AddFromFiles("mySettings", "www/base.html", "www/header.html", "www/accountSettings.html")
	r.AddFromFiles("projectAdd", "www/base.html", "www/header.html", "www/projectAdd.html")
	r.AddFromFiles("projectTaskBoard", "www/base.html", "www/header.html", "www/taskBoard.html").Funcs(InitRenderFunc())
	r.AddFromFiles("accountAdd", "www/base.html", "www/header.html", "www/accountAdd.html")

	return r
}

func InitRenderFunc() template.FuncMap {
	return template.FuncMap{
		"UnsafeHtml": func(text string) template.HTML { return template.HTML(text) },
	}
}
