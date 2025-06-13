package controller

import (
	"botdiscord/helper"

	"github.com/gin-gonic/gin"
)

// Importa o RenderWithUser do main
func Index(c *gin.Context) {
	helper.RenderWithUser(c, 200, "index.html", nil)
}
