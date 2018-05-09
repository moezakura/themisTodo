package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountView struct {
	*BaseView
}

func (self AccountView) GetAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "accountAdd", gin.H{
		"Title": "New Account",
	})
}
