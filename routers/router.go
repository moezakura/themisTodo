package routers

import (
	themsController "../controller"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB) *gin.Engine {
	r := gin.New()

	baseController := themsController.NewBaseController(db, r)

	// ログイン関連
	r.POST("/login", themsController.LoginController{baseController}.PostLogin)

	// プロジェクト関連
	projects := r.Group("/project")
	{
		projects.POST("/add", themsController.ProjectsController{baseController}.PostAdd)
		projects.POST("/delete/:projectId", themsController.ProjectsController{baseController}.PostDeleteProject)

		projects.POST("/update/:projectId", themsController.ProjectsController{baseController}.PostUpdate)
		projects.POST("/addUser/:projectId", themsController.ProjectsController{baseController}.PostAddUser)
		projects.GET("/info/:projectId", themsController.ProjectsController{baseController}.GetInfo)
		projects.GET("/tasks/:projectId", themsController.ProjectsController{baseController}.GetTasks)
		projects.GET("/members/:projectId", themsController.ProjectsController{baseController}.GetMembers)
		projects.DELETE("/members/:projectId", themsController.ProjectsController{baseController}.DeleteMembers)
		projects.GET("/my", themsController.ProjectsController{baseController}.GetMy)
	}

	// タスク管理
	tasks := r.Group("/tasks")
	{
		tasks.POST("/create", themsController.TasksController{baseController}.PostTaskCreate)
		tasks.POST("/update/:createDate", themsController.TasksController{baseController}.PostUpdate)
		tasks.POST("/delete/:createDate", themsController.TasksController{baseController}.PostDelete)
		tasks.GET("/view/:createDate", themsController.TasksController{baseController}.GetView)
		tasks.GET("/search", themsController.TasksController{baseController}.GetSearch)
		tasks.GET("/my", themsController.TasksController{baseController}.GetMy)
	}

	//アカウント関連
	account := r.Group("/account")
	{
		account.POST("/add", themsController.AccountController{baseController}.PostAdd)
		account.GET("/search", themsController.AccountController{baseController}.GetSearch)
		account.POST("/update", themsController.AccountController{baseController}.PostUpdate)
		account.POST("/updateIcon", themsController.AccountController{baseController}.PostUpdateIcon)
		account.GET("/profile", themsController.AccountController{baseController}.GetProfile)

		// icon
		account.GET("/icon/:iconPath", themsController.AccountController{baseController}.GetIcon)
	}
	return r
}