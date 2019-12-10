package routers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	themsController "themis.mox.si/themis/controller"
	"themis.mox.si/themis/module"
	"themis.mox.si/themis/repository"
)

func Init(db *sql.DB, gdb *gorm.DB, taskTimerWatcher *module.TaskTimerWatcherModule) *gin.Engine {
	r := gin.New()

	projectRepo := repository.NewProjectRepository(gdb)

	baseController := themsController.NewBaseController(db, r)

	// ログイン関連
	loginController := &themsController.LoginController{baseController}
	r.POST("/login", loginController.PostLogin)
	r.OPTIONS("/auth", loginController.AuthCheck)

	// プロジェクト関連
	projects := r.Group("/project")
	projects.Use(authCheck(db))
	{
		projectsController := themsController.NewProjectsController(baseController, projectRepo)

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
	tasks.Use(authCheck(db))
	{
		tasksController := themsController.NewTasksController(baseController, projectRepo)

		tasks.POST("/create", tasksController.PostTaskCreate)
		tasks.POST("/update/:createDate", tasksController.PostUpdate)
		tasks.POST("/bulkUpdate", tasksController.PostBulkUpdate)
		tasks.POST("/delete/:createDate", tasksController.PostDelete)
		tasks.DELETE("/bulkDelete", tasksController.DeleteBulkDelete)
		tasks.GET("/view/:createDate", tasksController.GetView)
		tasks.GET("/history/:createDate", tasksController.GetHistoryList)
		tasks.POST("/history/:createDate", tasksController.PostApplyHistory)
		tasks.GET("/search", tasksController.GetSearch)
		tasks.GET("/my", tasksController.GetMy)

		taskTimer := tasks.Group("/timer")
		{
			taskTimerController := themsController.NewTaskTimerController(baseController, taskTimerWatcher, projectRepo)
			taskTimer.PATCH("/toggle/:createDate", taskTimerController.PatchToggle)
			taskTimer.GET("/view/:createDate", taskTimerController.GetView)
			taskTimer.GET("/myDoing", taskTimerController.GetByUser)
			taskTimer.GET("/status/:createDate", taskTimerController.GetStatus)
			taskTimer.GET("/myList/:projectId", taskTimerController.GetMyList)
			taskTimer.DELETE("/delete/:taskTimerId", taskTimerController.Delete)
			taskTimer.POST("/update/:taskTimerId", taskTimerController.Update)
		}
	}

	//アカウント関連
	account := r.Group("/account")
	account.Use(authCheck(db))
	{
		accountsController := &themsController.AccountController{baseController}

		account.POST("/add", accountsController.PostAdd)
		account.GET("/search", accountsController.GetSearch)
		account.POST("/update", accountsController.PostUpdate)
		account.POST("/updateIcon", accountsController.PostUpdateIcon)
		account.GET("/profile", accountsController.GetProfile)
		account.GET("/list", accountsController.GetList)

		// icon
		r.GET("/account/icon/:iconPath", accountsController.GetIcon)
	}
	return r
}

func authCheck(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := map[string]interface{}{}
		token := c.GetHeader("X-Access-Token")

		session := module.NewSessionModule()
		exist, uuid := session.GetUuid(token)
		if !exist || len(token) <= 0 {
			res["message"] = "unknown token"
			res["success"] = false

			c.JSON(http.StatusForbidden, res)
			c.Abort()
			return
		} else {
			c.Set("uuid", uuid)
			c.Next()
		}
	}
}
