package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../models"
	"../module"
	"net/http"
	"strconv"
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
}

func (self ProjectsController) GetTaskBoard(c *gin.Context) {
	projectsModule := module.NewProjectsModule(self.DB)

	loginModule := module.NewLoginModule(self.DB)
	isError, accountUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	accountModule := module.NewAccountModule(self.DB)
	isError, account := accountModule.GetAccount(accountUuid)
	if isError {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	isError, project := projectsModule.GetProject(projectId)
	taskModule := module.NewTaskModule(self.DB)

	isErr, taskList := taskModule.GetList(projectId)
	if isErr {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	isErr, accounts := projectsModule.GetUser(projectId)
	if isErr {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	themisView.ProjectsView{self.BaseView}.GetTaskBoard(c, project, taskList, accounts, account)
}

func (self ProjectsController) PostUpdate(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}
	loginModule := module.NewLoginModule(self.DB)
	isError, _ := loginModule.GetUserId(c, self.Session)

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

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

	//TODO: ユーザーがProjectに参加してるかチェック

	project := &models.Project{
		projectId,
		addRequest.Name,
		addRequest.Description,
	}

	projectsModule := module.NewProjectsModule(self.DB)
	err2 := projectsModule.Update(project)

	if err2 {
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	addResult.Success = true
	themisView.ProjectsView{}.PostUpdate(c, addResult)
}

func (self ProjectsController) PostAddUser(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}
	loginModule := module.NewLoginModule(self.DB)
	isError, _ := loginModule.GetUserId(c, self.Session)

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	if isError {
		addResult.Message = "invalid token"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
		return
	}

	var addRequest models.ProjectsAddUserRequest
	c.ShouldBindJSON(&addRequest)

	if addRequest.Uuid < 0 {
		addResult.Message = "invalid user id"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
		return
	}

	accountModule := module.NewAccountModule(self.DB)

	isError, account := accountModule.GetAccount(addRequest.Uuid)
	if isError || account == nil {
		addResult.Message = "not found user id"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
		return
	}

	searchObject := models.NewAccountSearchModel()
	searchObject.ProjectId = projectId
	searchObject.IsInProject = true
	searchObject.Uuid = addRequest.Uuid

	isError, searchResult := accountModule.Search(searchObject)
	if isError {
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
		return
	}

	if len(searchResult) > 0 {
		addResult.Message = "this user has already joined the project"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
		return
	}

	projectModule := module.NewProjectsModule(self.DB)
	isErrorProjectAdd := projectModule.AddUser(addRequest.Uuid, projectId)
	if isErrorProjectAdd {
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
	}

	addResult.AddedAccount = &models.Account{
		Name:        account.Name,
		DisplayName: account.DisplayName,
		Uuid:        account.Uuid,
	}

	addResult.Success = true
	themisView.ProjectsView{}.PostAddUser(c, addResult)
}

func (self ProjectsController) PostDeleteProject(c *gin.Context) {
	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)
	resultJson := models.ProjectDeleteResultJson{}
	if isError {
		resultJson.Message = "invalid token"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusBadRequest, &resultJson)
		return
	}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusBadRequest, &resultJson)
		return
	}

	projectModule := module.NewProjectsModule(self.DB)
	isIn := projectModule.IsIn(userUuid, projectId)
	if !isIn {
		resultJson.Message = "invalid user"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusBadRequest, &resultJson)
		return
	}

	isError = projectModule.Delete(projectId)
	if isError {
		resultJson.Message = "failed delete"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusInternalServerError, &resultJson)
		return
	}

	resultJson.Message = "ok"
	resultJson.Success = true
	themisView.ProjectsView{}.PostDeleteProject(c, http.StatusOK, &resultJson)
}

func (self ProjectsController) GetMy(c *gin.Context) {
	getResult := &models.ProjectGetResultJson{}
	projectsModule := module.NewProjectsModule(self.DB)

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		getResult.Message = "invalid token"
		themisView.ProjectsView{}.GetMy(c, http.StatusBadRequest, getResult)
		return
	}

	isError, projects := projectsModule.GetProjects(userUuid)
	if isError {
		getResult.Message = "unknown project"
		themisView.ProjectsView{}.GetMy(c, http.StatusBadRequest, getResult)
		return
	}

	getResult.Project = projects
	getResult.Success = true

	themisView.ProjectsView{}.GetMy(c, http.StatusOK, getResult)

}
