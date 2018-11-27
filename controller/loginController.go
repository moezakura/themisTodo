package controller

import (
	"../models"
	"../module"
	"../utils"
	themisView "../view"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
	*BaseController
}

func (l *LoginController) PostLogin(c *gin.Context) {
	var loginRequest models.LoginRequestJson
	c.ShouldBindJSON(&loginRequest)

	loginResult := &models.LoginResultJson{}
	if len(loginRequest.Password) < 1 || len(loginRequest.Id) < 1 {
		loginResult.Message = "id and password is not allowed empty"
		themisView.LoginView{}.PostLogin(c, loginResult)
		return
	}

	loginRequest.Password = utils.SHA512(loginRequest.Password)
	loginModule := module.NewLoginModule(l.DB)

	err, uuid := loginModule.IsLogin(loginRequest.Id, loginRequest.Password)

	if err {
		loginResult.Message = "invalid id or password"
		themisView.LoginView{}.PostLogin(c, loginResult)
		return
	}

	authToken := l.Session.GetToken(uuid)

	loginResult.Success = true
	loginResult.Message = authToken

	c.JSON(http.StatusOK, loginResult)
}

func (l *LoginController) AuthCheck(c *gin.Context) {
	loginModule := module.NewLoginModule(l.DB)
	isError, _ := loginModule.GetUserId(c, l.Session)

	loginResult := &models.LoginResultJson{}
	status := http.StatusBadRequest

	if isError {
		loginResult.Success = false
		loginResult.Message = "invalid token"
	} else {
		loginResult.Success = true
		loginResult.Message = "ok"
		status = http.StatusOK
	}

	c.JSON(status, loginResult)
}
