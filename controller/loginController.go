package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../models"
	"../utils"
	"../module"
	"net/http"
)

type LoginController struct {
	*BaseController
}

func (self LoginController) GetLogin(c *gin.Context) {
	themisView.LoginView{self.BaseView}.GetLogin(c)
}

func (self LoginController) PostLogin(c *gin.Context) {
	var loginRequest models.LoginRequestJson
	c.ShouldBindJSON(&loginRequest)

	loginResult := &models.LoginResultJson{}
	if len(loginRequest.Password) < 1 || len(loginRequest.Id) < 1 {
		loginResult.Message = "id and password is not allowed empty"
		themisView.LoginView{}.PostLogin(c, loginResult)
		return
	}

	loginRequest.Password = utils.SHA512(loginRequest.Password)
	loginModule := module.NewLoginModule(self.DB)

	err, uuid := loginModule.IsLogin(loginRequest.Id, loginRequest.Password)

	if err {
		loginResult.Message = "invalid id or password"
		themisView.LoginView{}.PostLogin(c, loginResult)
		return
	}

	limitSec := 30 * 24 * 60 * 60
	authToken := utils.RandomString(64)
	self.Session.Add(authToken, uuid, limitSec)

	c.SetCookie("token", authToken, limitSec, "/", "", false, true)
	loginResult.Success = true
	loginResult.Message = authToken

	c.JSON(http.StatusOK, loginResult)
}
