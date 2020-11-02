package main

import (
	"github.com/epicadk/Lunchbox-backend/src/handler"
	"github.com/epicadk/Lunchbox-backend/src/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/login", handler.LoginPost)
	r.POST("/signup", handler.SignupPost)
	r.POST("/refresh", handler.RefreshPost)

	//setting up middleware
	r.Use(middleware.TokenAuthMiddleware)

	r.POST("/validate", handler.AuthTokenPost)
	r.POST("/comments", handler.CommentsPost)
	r.GET("/comments", handler.CommentsGet)
	r.POST("/fav", handler.FavRestaurantPost)
	r.GET("/fav", handler.FavouriteRestaurantGet)
	r.Run() //default is 8080
}
