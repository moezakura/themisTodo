package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	themisView "../view"
	"../models"
	"../module"
	"../utils"
	"net/http"
	"strconv"
	"math"
	"fmt"
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

	loginModule := module.NewLoginModule(self.DB)

	isErr, accountUuid := loginModule.GetUserId(c, self.Session)
	if isErr {
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
	var updateRequest models.AccountChangeRequestJson
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

	accountChangeRequest := &models.AccountChangeRequestJson{}

	if len(updateRequest.Password) > 0 {
		if len(updateRequest.Password) < 10 {
			result.Message = "password must be at least 10 letters"
			themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
			return
		}

		passwordHash := utils.SHA512(updateRequest.CurrentPassword)
		isErr, _ := loginModule.IsLoginFromUuid(accountUuid, passwordHash)
		if !isErr {
			accountChangeRequest.Password = updateRequest.Password
			isChange = true
		} else {
			result.Message = "the current password is wrong"
			themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
			return
		}
	}

	accountChangeRequest.DisplayName = account.DisplayName
	accountChangeRequest.Name = account.Name
	accountChangeRequest.Uuid = account.Uuid

	if !isChange {
		result.Message = "no change"
		themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusBadRequest, &result)
		return
	}

	updateRequest.Uuid = accountUuid
	accountModule.Update(accountChangeRequest)
	result.Success = true
	themisView.AccountView{self.BaseView}.PostUpdate(c, http.StatusOK, &result)
}

func (self AccountController) PostUpdateIcon(c *gin.Context) {
	result := models.AccountUpdateIconRequest{}
	iconFile, err := c.FormFile("icon")
	if err != nil {
		result.Message = "icon file is not found"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusBadRequest, &result)
		return
	}

	if iconFile.Size > 3*1024*1024 {
		result.Message = "maximum file size is 3MB"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusBadRequest, &result)
		return
	}

	accountUuidStr := c.Param("accountUuid")
	accountUuidI64, err := strconv.ParseInt(accountUuidStr, 10, 32)
	accountUuid := int(accountUuidI64)

	if err != nil {
		result.Message = "invalid account id"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusBadRequest, &result)
		return
	}

	loginModule := module.NewLoginModule(self.DB)

	isErr, sessionUuid := loginModule.GetUserId(c, self.Session)
	if sessionUuid != accountUuid || isErr {
		result.Message = "invalid account id"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusBadRequest, &result)
		return
	}

	iconSavePath := fmt.Sprintf("www/assets/accountIcon/%d.png", accountUuid)

	img, err := imageupload.Process(c.Request, "icon")
	if err != nil {
		result.Message = "invalid image file"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusBadRequest, &result)
		return
	}
	thumb, err := imageupload.ThumbnailPNG(img, 500, 500)
	if err != nil {
		result.Message = "invalid image file"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusBadRequest, &result)
		return
	}
	thumb.Save(iconSavePath)

	result.Success = true
	themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusOK, &result)
}
