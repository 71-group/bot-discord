package controller

import "github.com/gin-gonic/gin"

func GetUserList(c *gin.Context) {
	context := gin.H{
		"User": "MTzica",
	}
	c.HTML(200, "user-list.html", context)
}
