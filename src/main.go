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

	// Comments endpoints
	r.POST("/comments", handler.CommentsPost)
	r.DELETE("/comments", handler.CommentsDelete)
	r.PUT("/comments", handler.CommentsPut)
	r.GET("/restaurant/comments", handler.CommentsGet)
	r.GET("/user/comments", handler.RecentActivityGet)

	// User Favourite Restaurants endpoints
	r.POST("/fav", handler.FavRestaurantPost)
	r.GET("/fav", handler.FavouriteRestaurantGet)

	_ = r.Run() //default is 8080
}
