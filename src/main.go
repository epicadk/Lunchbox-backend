package main

import (
	handler "github.com/epicadk/Lunchbox-backend/src/Handler"
	middleware "github.com/epicadk/Lunchbox-backend/src/Middleware"

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
	r.GET("/fav", handler.FavouriteRestaurantGet)
	r.Run() //default is 8080
}
