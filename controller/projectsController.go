package controller

import (
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/module"
	"themis.mox.si/themis/utils"
	themisView "themis.mox.si/themis/view"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProjectsController struct {
	*BaseController
}

func (p *ProjectsController) PostAdd(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}
	userUuid := c.GetInt("uuid")

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

	projectsModule := module.NewProjectsModule(p.DB)
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

func (p *ProjectsController) PostUpdate(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request")
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

	projectsModule := module.NewProjectsModule(p.DB)
	err2 := projectsModule.Update(project)

	if err2 {
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	addResult.Success = true
	themisView.ProjectsView{}.PostUpdate(c, addResult)
}

func (p *ProjectsController) PostAddUser(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	var addRequest models.ProjectsAddUserRequest
	c.ShouldBindJSON(&addRequest)

	if addRequest.Uuid < 0 {
		addResult.Message = "invalid user id"
		themisView.ProjectsView{}.PostAddUser(c, addResult)
		return
	}

	accountModule := module.NewAccountModule(p.DB)

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

	projectModule := module.NewProjectsModule(p.DB)
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

func (p *ProjectsController) PostDeleteProject(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	resultJson := models.ProjectDeleteResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusBadRequest, &resultJson)
		return
	}

	projectModule := module.NewProjectsModule(p.DB)
	isIn := projectModule.IsIn(userUuid, projectId)
	if !isIn {
		resultJson.Message = "invalid user"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusBadRequest, &resultJson)
		return
	}

	isError := projectModule.Delete(projectId)
	if isError {
		resultJson.Message = "failed delete"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusInternalServerError, &resultJson)
		return
	}

	resultJson.Message = "ok"
	resultJson.Success = true
	themisView.ProjectsView{}.PostDeleteProject(c, http.StatusOK, &resultJson)
}

func (p *ProjectsController) GetInfo(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	resultJson := &models.ProjectInfoResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.GetInfo(c, http.StatusBadRequest, resultJson)
		return
	}

	projectModule := module.NewProjectsModule(p.DB)
	isIn := projectModule.IsIn(userUuid, projectId)
	if !isIn {
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetInfo(c, http.StatusNotFound, resultJson)
		return
	}

	isError, project := projectModule.GetProject(projectId)
	if isError {
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetInfo(c, http.StatusNotFound, resultJson)
		return
	}
	resultJson.Project = *project
	resultJson.Success = true

	themisView.ProjectsView{}.GetInfo(c, http.StatusOK, resultJson)
}

func (p *ProjectsController) GetTasks(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	resultJson := &models.TasksGetResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.GetTasks(c, http.StatusBadRequest, resultJson)
		return
	}

	projectModule := module.NewProjectsModule(p.DB)
	isIn := projectModule.IsIn(userUuid, projectId)
	if !isIn {
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetTasks(c, http.StatusNotFound, resultJson)
		return
	}

	isError, _ := projectModule.GetProject(projectId)
	if isError {
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetTasks(c, http.StatusNotFound, resultJson)
		return
	}

	tasksModule := module.NewTaskModule(p.DB)
	isError, tasks := tasksModule.GetList(projectId)
	tasks = utils.TasksConvert(tasks)

	resultJson.Success = true
	resultJson.Task = models.NewTasksOfJson(tasks)

	themisView.ProjectsView{}.GetTasks(c, http.StatusOK, resultJson)
}

func (p *ProjectsController) GetMembers(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	resultJson := &models.ProjectMembersResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.GetMembers(c, http.StatusBadRequest, resultJson)
		return
	}

	projectModule := module.NewProjectsModule(p.DB)
	isIn := projectModule.IsIn(userUuid, projectId)
	if !isIn {
		resultJson.Message = "not found project or not found users"
		themisView.ProjectsView{}.GetMembers(c, http.StatusNotFound, resultJson)
		return
	}

	isError, users := projectModule.GetUser(projectId)
	if isError {
		resultJson.Message = "not found project or not found users"
		themisView.ProjectsView{}.GetMembers(c, http.StatusNotFound, resultJson)
		return
	}

	resultJson.Success = true
	resultJson.Members = users

	themisView.ProjectsView{}.GetMembers(c, http.StatusOK, resultJson)
}

func (p *ProjectsController) DeleteMembers(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	resultJson := &models.ProjectDeleteMemberResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusBadRequest, resultJson)
		return
	}

	projectModule := module.NewProjectsModule(p.DB)
	isIn := projectModule.IsIn(userUuid, projectId)
	if !isIn {
		resultJson.Message = "not found project or not found users"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusNotFound, resultJson)
		return
	}

	var deleteRequest models.ProjectsDeleteUserRequest
	c.ShouldBindJSON(&deleteRequest)

	if deleteRequest.Uuid < 0 {
		resultJson.Message = "invalid user id"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusBadRequest, resultJson)
		return
	}

	isInTarget := projectModule.IsIn(deleteRequest.Uuid, projectId)
	if !isInTarget {
		resultJson.Message = "not found project or not found users"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusNotFound, resultJson)
		return
	}

	isError := projectModule.Leave(projectId, deleteRequest.Uuid)
	if isError {
		resultJson.Message = "failed leave"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusInternalServerError, resultJson)
		return
	}

	resultJson.Success = true

	themisView.ProjectsView{}.DeleteMember(c, http.StatusOK, resultJson)
}

func (p *ProjectsController) GetMy(c *gin.Context) {
	getResult := &models.ProjectGetResultJson{}
	projectsModule := module.NewProjectsModule(p.DB)

	userUuid := c.GetInt("uuid")

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
