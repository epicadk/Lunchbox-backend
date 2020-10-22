package handler

import (
	"github.com/gin-gonic/gin"
)

//AuthTokenPost function handles post request to authtoken endpoint
func AuthTokenPost(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "valid user",
	})
}
