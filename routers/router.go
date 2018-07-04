package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/multitemplate"
	themsController "../controller"
	"database/sql"
	"html/template"
	"github.com/jinzhu/gorm"
)

func Init(db *sql.DB, gormDB *gorm.DB) *gin.Engine {
	r := gin.New()

	r.Static("/fontawesome", "./www/assets/fontawesome/web-fonts-with-css/")
	r.Static("/assets", "./www/assets")
	r.SetFuncMap(InitRenderFunc())
	r.HTMLRender = InitRender()

	gormDB.LogMode(true)
	baseController := themsController.NewBaseController(db, gormDB, r)

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
		tasks.POST("/delete/:createDate", themsController.TasksController{baseController}.PostDelete)
		tasks.GET("/view/:createDate", themsController.TasksController{baseController}.GetView)
		tasks.GET("/search", themsController.TasksController{baseController}.GetSearch)
		tasks.GET("/searches", themsController.TasksController{baseController}.GetSearches)
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
	funcMap := InitRenderFunc()

	r.AddFromFilesFuncs("index", funcMap, "www/base.html", "www/index.html").Funcs(funcMap)
	r.AddFromFilesFuncs("login", funcMap, "www/base.html", "www/login.html").Funcs(funcMap)
	r.AddFromFilesFuncs("home", funcMap, "www/base.html", "www/header.html", "www/home.html").Funcs(funcMap)
	r.AddFromFilesFuncs("mySettings", funcMap, "www/base.html", "www/header.html", "www/accountSettings.html").Funcs(funcMap)
	r.AddFromFilesFuncs("projectAdd", funcMap, "www/base.html", "www/header.html", "www/projectAdd.html").Funcs(funcMap)
	r.AddFromFilesFuncs("projectTaskBoard", funcMap, "www/base.html", "www/header.html", "www/taskBoard.html").Funcs(funcMap)
	r.AddFromFilesFuncs("accountAdd", funcMap, "www/base.html", "www/header.html", "www/accountAdd.html").Funcs(funcMap)

	return r
}

func InitRenderFunc() template.FuncMap {
	return template.FuncMap{
		"UnsafeHtml": func(text string) template.HTML { return template.HTML(text) },
		"isDebug":    func() bool { return gin.Mode() == gin.DebugMode },
	}
}
