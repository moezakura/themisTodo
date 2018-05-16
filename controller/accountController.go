package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../models"
	"../module"
	"../utils"
	"net/http"
	"strconv"
	"math"
)

type AccountController struct {
	*BaseController
}

func (self AccountController) GetAdd(c *gin.Context) {
	themisView.AccountView{self.BaseView}.GetAdd(c)
}

func (self AccountController) PostAdd(c *gin.Context) {
	var addRequest models.AccountAddRequestJson
	c.ShouldBindJSON(&addRequest)

	addResult := &models.AccountAddResultJson{}

	if len(addRequest.Name) < 1 {
		addResult.Message = "id is not allowed empty"
		themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
		return
	}

	if len(addRequest.Name) > 128 {
		addResult.Message = "maximum name length is 128 characters"
		themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
		return
	}

	accountModule := module.NewAccountModule(self.DB)

	if accountModule.Get(addRequest.Name) > 0 {
		addResult.Message = "this id is exist"
		themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
		return
	}

	// TODO: ロールによる権限管理の追加

	password := utils.RandomString(24)

	accountModule.Add(addRequest.Name, password)

	addResult.Name = addRequest.Name
	addResult.Success = true
	addResult.Password = password
	themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
}

func (self AccountController) GetSearch(c *gin.Context) {
	projectIdTmp := c.DefaultQuery("project", "")
	isInProjectTmp := c.DefaultQuery("isInProject", "")
	displayNameTemp := c.DefaultQuery("displayName", "")
	nameTemp := c.DefaultQuery("name", "")
	uuidTemp := c.DefaultQuery("uuid", "")
	maxTemp := c.DefaultQuery("max", "")

	searchModel := models.NewAccountSearchModel()
	if projectIdTmp != "" {
		projectIdNum, err := strconv.ParseInt(projectIdTmp, 10, 64)
		isInProjectBool, errB := strconv.ParseBool(isInProjectTmp)
		if errB == nil {
			searchModel.IsInProject = isInProjectBool
		}
		if err == nil && projectIdNum < math.MaxInt32 {
			searchModel.ProjectId = int(projectIdNum)
		}
	}
	if displayNameTemp != "" {
		searchModel.DisplayName = displayNameTemp
	}
	if nameTemp != "" {
		searchModel.Name = nameTemp
	}
	if maxTemp != "" {
		maxNum, err := strconv.ParseInt(maxTemp, 10, 64)
		if err == nil && maxNum < math.MaxInt32 {
			searchModel.Max = int(maxNum)
		}
	}
	if uuidTemp != "" {
		uuidNum, err := strconv.ParseInt(uuidTemp, 10, 64)
		if err == nil && uuidNum < math.MaxInt32 {
			searchModel.Uuid = int(uuidNum)
		}
	}

	accountModule := module.NewAccountModule(self.DB)
	isError, result := accountModule.Search(searchModel)
	if !isError {
		themisView.AccountView{self.BaseView}.GetSearch(c, http.StatusOK, &result)
	} else {
		themisView.AccountView{self.BaseView}.GetSearch(c, http.StatusBadRequest, nil)
	}
}

func (self AccountController) PostUpdate(c *gin.Context) {
	result := models.AccountChangeResultJson{}

	accountUuidStr := c.Param("accountUuid")
	accountUuidI64, err := strconv.ParseInt(accountUuidStr, 10, 32)
	accountUuid := int(accountUuidI64)

	if err != nil {
		result.Message = "invalid account id"
		themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
		return
	}

	accountModule := module.NewAccountModule(self.DB)
	isErr, account := accountModule.GetAccount(accountUuid)

	if isErr {
		result.Message = "invalid account id"
		themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
		return
	}

	isChange := false
	var updateRequest models.Account
	c.ShouldBindJSON(&updateRequest)

	if len(updateRequest.Name) > 128 {
		result.Message = "maximum id length is 128 characters"
		themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
		return
	}

	if len(updateRequest.DisplayName) > 128 {
		result.Message = "maximum display name length is 128 characters"
		themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
		return
	}

	if len(updateRequest.DisplayName) > 0 {
		account.DisplayName = updateRequest.DisplayName
		isChange = true
	}
	if len(updateRequest.Name) > 0 {
		account.Name = updateRequest.Name
		isChange = true
	}
	if len(updateRequest.Password) > 0 {
		account.Password = updateRequest.Password
		isChange = true
	}

	if !isChange{
		result.Message = "no change"
		themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
		return
	}

	updateRequest.Uuid = accountUuid
	accountModule.Update(account)
	result.Success = true
	themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
}
