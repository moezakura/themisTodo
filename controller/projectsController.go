package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/models/db"
	"themis.mox.si/themis/models/json"
	"themis.mox.si/themis/module"
	"themis.mox.si/themis/repository"
	"themis.mox.si/themis/utils"
	themisView "themis.mox.si/themis/view"
)

type ProjectsController struct {
	*BaseController
	projectRepo *repository.ProjectRepository
	taskRepo    *repository.TaskRepository
}

func NewProjectsController(baseController *BaseController, projectRepo *repository.ProjectRepository, taskRepo *repository.TaskRepository) *ProjectsController {
	return &ProjectsController{
		BaseController: baseController,
		projectRepo:    projectRepo,
		taskRepo:       taskRepo,
	}
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

	id, err := p.projectRepo.Add(addRequest.Name, addRequest.Description)

	if err != nil {
		// TODO: Logにエラーを残す
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}
	err = p.projectRepo.AddUser(userUuid, id)
	if err != nil {
		// TODO: Logにエラーを残す
		// TODO: エラーを返す
	}

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

	project := &db.Project{
		Uuid:        projectId,
		Name:        addRequest.Name,
		Description: addRequest.Description,
	}
	err = p.projectRepo.Update(project)

	if err != nil {
		// TODO: Logにエラーを残す
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

	err = p.projectRepo.AddUser(addRequest.Uuid, projectId)
	if err != nil {
		// TODO: Logにエラーを残す
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

	isIn, err := p.projectRepo.IsIn(userUuid, projectId)
	// TODO: Logにエラーを残す
	// TODO: エラーを返す
	if !isIn {
		resultJson.Message = "invalid user"
		themisView.ProjectsView{}.PostDeleteProject(c, http.StatusBadRequest, &resultJson)
		return
	}

	err = p.projectRepo.Delete(projectId)
	if err != nil {
		// TODO: Logにエラーを残す
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
	resultJson := &json.ProjectInfoResultJson{}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		resultJson.Message = "invalid project id"
		themisView.ProjectsView{}.GetInfo(c, http.StatusBadRequest, resultJson)
		return
	}

	isIn, err := p.projectRepo.IsIn(userUuid, projectId)
	// TODO: Logにエラーを残す
	// TODO: エラーを返す
	if !isIn {
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetInfo(c, http.StatusNotFound, resultJson)
		return
	}

	project, err := p.projectRepo.GetProjectById(projectId)
	if err != nil {
		// TODO: Logにエラーを残す
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

	isIn, err := p.projectRepo.IsIn(userUuid, projectId)
	// TODO: Logにエラーを残す
	// TODO: エラーを返す
	if !isIn {
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetTasks(c, http.StatusNotFound, resultJson)
		return
	}

	_, err = p.projectRepo.GetProjectById(projectId)
	if err != nil {
		// TODO: Logにエラーを残す
		resultJson.Message = "not found project"
		themisView.ProjectsView{}.GetTasks(c, http.StatusNotFound, resultJson)
		return
	}

	tasks, err := p.taskRepo.GetList(projectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
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

	isIn, err := p.projectRepo.IsIn(userUuid, projectId)
	// TODO: Logにエラーを残す
	// TODO: エラーを返す
	if !isIn {
		resultJson.Message = "not found project or not found users"
		themisView.ProjectsView{}.GetMembers(c, http.StatusNotFound, resultJson)
		return
	}

	users, err := p.projectRepo.GetUserListInProject(projectId)
	if err != nil {
		// TODO: Logにエラーを残す
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

	isIn, err := p.projectRepo.IsIn(userUuid, projectId)
	// TODO: Logにエラーを残す
	// TODO: エラーを返す
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

	isInTarget, err := p.projectRepo.IsIn(deleteRequest.Uuid, projectId)
	// TODO: Logにエラーを残す
	// TODO: エラーを返す
	if !isInTarget {
		resultJson.Message = "not found project or not found users"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusNotFound, resultJson)
		return
	}

	err = p.projectRepo.Leave(projectId, deleteRequest.Uuid)
	if err != nil {
		// TODO: Logにエラーを残す
		resultJson.Message = "failed leave"
		themisView.ProjectsView{}.DeleteMember(c, http.StatusInternalServerError, resultJson)
		return
	}

	resultJson.Success = true

	themisView.ProjectsView{}.DeleteMember(c, http.StatusOK, resultJson)
}

func (p *ProjectsController) GetMy(c *gin.Context) {
	getResult := &json.ProjectGetResultJson{}

	userUuid := c.GetInt("uuid")

	projects, err := p.projectRepo.GetProjectsByUserId(userUuid)
	if err != nil {
		// TODO: Logにエラーを残す
		getResult.Message = "unknown project"
		themisView.ProjectsView{}.GetMy(c, http.StatusBadRequest, getResult)
		return
	}

	getResult.Project = projects
	getResult.Success = true

	themisView.ProjectsView{}.GetMy(c, http.StatusOK, getResult)

}
