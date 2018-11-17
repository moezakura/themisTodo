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
		projectsController := themsController.ProjectsController{baseController}

		projects.POST("/add", projectsController.PostAdd)
		projects.POST("/delete/:projectId", projectsController.PostDeleteProject)

		projects.POST("/update/:projectId", projectsController.PostUpdate)
		projects.POST("/addUser/:projectId", projectsController.PostAddUser)
		projects.GET("/info/:projectId", projectsController.GetInfo)
		projects.GET("/tasks/:projectId", projectsController.GetTasks)
		projects.GET("/members/:projectId", projectsController.GetMembers)
		projects.DELETE("/members/:projectId", projectsController.DeleteMembers)
		projects.GET("/my", projectsController.GetMy)
	}

	// タスク管理
	tasks := r.Group("/tasks")
	{
		tasksController := themsController.TasksController{baseController}

		tasks.POST("/create", tasksController.PostTaskCreate)
		tasks.POST("/update/:createDate", tasksController.PostUpdate)
		tasks.POST("/delete/:createDate", tasksController.PostDelete)
		tasks.GET("/view/:createDate", tasksController.GetView)
		tasks.GET("/search", tasksController.GetSearch)
		tasks.GET("/my", tasksController.GetMy)
	}

	//アカウント関連
	account := r.Group("/account")
	{
		accountsController := themsController.AccountController{baseController}

		account.POST("/add", accountsController.PostAdd)
		account.GET("/search", accountsController.GetSearch)
		account.POST("/update", accountsController.PostUpdate)
		account.POST("/updateIcon", accountsController.PostUpdateIcon)
		account.GET("/profile", accountsController.GetProfile)
		account.GET("/list", accountsController.GetList)

		// icon
		account.GET("/icon/:iconPath", accountsController.GetIcon)
	}
	return r
}