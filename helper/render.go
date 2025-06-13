package helper

import (
	"github.com/gin-gonic/gin"
)

func RenderWithUser(c *gin.Context, code int, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	data["UserID"] = c.MustGet("UserID")
	data["Username"] = c.MustGet("Username")
	data["Avatar"] = c.MustGet("Avatar")
	c.HTML(code, name, data)
}
