package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"themis.mox.si/themis/models"
)

type AccountView struct {
	*BaseView
}

func (self AccountView) PostAdd(c *gin.Context, json *models.AccountAddResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self AccountView) GetSearch(c *gin.Context, statusCode int, json *[]models.AccountSearchResultModel) {
	c.JSON(statusCode, json)
}

func (self AccountView) PostUpdate(c *gin.Context, statusCode int, json *models.AccountChangeResultJson)  {
	c.JSON(statusCode, json)
}

func (self AccountView) PostUpdateIcon(c *gin.Context, statusCode int, json *models.AccountUpdateIconRequest)  {
	c.JSON(statusCode, json)
}

func (self AccountView) GetProfile(c *gin.Context, statusCode int, json *models.AccountProfileResultJson)  {
	c.JSON(statusCode, json)
}

func (self AccountView) GetList(c *gin.Context, statusCode int, json *models.AccountListResult)  {
	c.JSON(statusCode, json)
}
