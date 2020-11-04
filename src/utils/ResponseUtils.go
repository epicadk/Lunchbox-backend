package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondWithQuickError(c *gin.Context, httpStatus int) {
	c.AbortWithStatusJSON(httpStatus, gin.H{
		"message": http.StatusText(httpStatus),
	})
}

func RespondWithError(c *gin.Context, httpStatus int, message string) {
	c.AbortWithStatusJSON(httpStatus, gin.H{
		"message": message,
	})
}
