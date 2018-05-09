package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../module"
	"net/http"
)

type HomeController struct {
	*BaseController
}

func (self HomeController) GetHome(c *gin.Context) {
	projectsModule := module.NewProjectsModule(self.DB)

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	isError, projects := projectsModule.GetProjects(userUuid)

	themisView.HomeView{self.BaseView}.GetHome(c, projects)
}
