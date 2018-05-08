package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../models"
	"../module"
)

type ProjectsController struct {
	*BaseController
}

func (self ProjectsController) GetAdd(c *gin.Context) {
	themisView.ProjectsView{self.BaseView}.GetAdd(c)
}

func (self ProjectsController) PostAdd(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}
	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		addResult.Message = "invalid token"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	var addRequest models.ProjectAddRequestJson
	c.ShouldBindJSON(&addRequest)

	if len(addRequest.Name) < 1 || len(addRequest.Description) < 1 {
		addResult.Message = "name and description is not allowed empty"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	if len(addRequest.Name) > 256 {
		addResult.Message = "maximum name length is 256 characters"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	if len(addRequest.Description) > 1024 {
		addResult.Message = "maximum description length is 1024 characters"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	projectsModule := module.NewProjectsModule(self.DB)
	err2, id := projectsModule.Add(addRequest.Name, addRequest.Description)

	if err2 {
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}
	projectsModule.AddUser(userUuid, id)

	addResult.Success = true
	addResult.Id = id
	themisView.ProjectsView{}.PostAdd(c, addResult)

	return
}
