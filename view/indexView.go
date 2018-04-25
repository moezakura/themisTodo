package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexView struct {
}

func (self IndexView) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
