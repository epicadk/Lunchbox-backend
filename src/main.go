package main

import (
	handler "go-gin-api/src/Handler"
	middleware "go-gin-api/src/Middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/login", handler.LoginPost)
	r.POST("/signup", handler.SignupPost)

	//setting up middleware
	r.Use(middleware.TokenAuthMiddleware)

	r.POST("/validate", handler.AuthTokenPost)
	r.POST("/comments", handler.CommentsPost)
	r.GET("/comments", handler.CommentsGet)
	r.POST("/fav", handler.FavRestaurantPost)
	r.Run() //default is 8080
}
