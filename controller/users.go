package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Usuario struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

func GetUserList(c *gin.Context) {
	usuarios := []Usuario{
		{ID: "123", Nome: "Matheus"},
		{ID: "456", Nome: "Vitor"},
	}
	c.JSON(http.StatusOK, usuarios)
}
