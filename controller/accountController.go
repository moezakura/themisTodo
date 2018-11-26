package controller

import (
	"../models"
	"../module"
	"../utils"
	themisView "../view"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type AccountController struct {
	*BaseController
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
	userUuid := c.GetInt("uuid")

	accountModule := module.NewAccountModule(self.DB)
	isErr, account := accountModule.GetAccount(userUuid)

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
		isErr, _ := loginModule.IsLoginFromUuid(userUuid, passwordHash)
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

	updateRequest.Uuid = userUuid
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

	userUuid := c.GetInt("uuid")

	accountModule := module.NewAccountModule(self.DB)
	isErr, account := accountModule.GetAccount(userUuid)
	if isErr {
		result.Message = "server error"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusInternalServerError, &result)
		return
	}
Retry:
	imageName := utils.RandomString(48)
	if accountModule.IsExistFromIconPath(imageName) {
		goto Retry
	}

	iconSavePath := fmt.Sprintf("data/account_icon/%s.png", imageName)

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
	err = thumb.Save(iconSavePath)
	if err != nil {
		result.Message = "image save error"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusInternalServerError, &result)
		return
	}

	oldImageName := account.IconPath
	account.IconPath = imageName
	isErr = accountModule.UpdateIconPath(userUuid, imageName)
	if isErr {
		os.Remove(iconSavePath)
		result.Message = "image save error"
		themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusInternalServerError, &result)
		return
	}

	oldImagePath := fmt.Sprintf("data/account_icon/%s.png", oldImageName)
	_, iErr := os.Stat(oldImagePath)
	if !os.IsNotExist(iErr) {
		err = os.Remove(oldImagePath)
		if err != nil {
			result.Message = "image save error"
			log.Printf("oldImage remove error: %+v", err)
			themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusInternalServerError, &result)
			return
		}
	}

	result.Success = true
	result.FileId = imageName
	themisView.AccountView{self.BaseView}.PostUpdateIcon(c, http.StatusOK, &result)
}

func (self AccountController) GetProfile(c *gin.Context) {
	profileResult := &models.AccountProfileResultJson{}
	userUuid := c.GetInt("uuid")

	accountModule := module.NewAccountModule(self.DB)
	isError, account := accountModule.GetAccount(userUuid)
	if isError {
		profileResult.Message = "server error"
		themisView.AccountView{}.GetProfile(c, http.StatusInternalServerError, profileResult)
		return
	}

	profileResult.Success = true
	profileResult.User = account
	themisView.AccountView{}.GetProfile(c, http.StatusOK, profileResult)
}

func (self AccountController) GetList(c *gin.Context) {
	result := &models.AccountListResult{}

	accountModule := module.NewAccountModule(self.DB)
	isError, account := accountModule.GetAccountsList()
	if isError {
		result.Message = "server error"
		themisView.AccountView{}.GetList(c, http.StatusInternalServerError, result)
		return
	}

	result.Success = true
	result.Users = account
	themisView.AccountView{}.GetList(c, http.StatusOK, result)
}

func (self AccountController) GetIcon(c *gin.Context) {
	iconPathQuery := c.Param("iconPath")
	iconPath := filepath.Join("data/account_icon/", iconPathQuery+".png")

	_, err := os.Stat(iconPath)
	if !os.IsNotExist(err) {
		c.File(iconPath)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/assets/images/unknown.png")
	}
}
