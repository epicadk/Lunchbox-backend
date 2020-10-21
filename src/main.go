package main

import (
	handler "go-gin-api/src/Handler"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/login", handler.LoginPost)
	r.POST("/signup", handler.SignupPost)
	r.Run() //default is 8080
}
