package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func FavouriteRestaurantGet(c *gin.Context) {
	//TODO change to url param
	username := c.Query("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := database.CheckUserFromUsername(ctx, c, username)
	if err != nil {
		return
	}
	client := database.MongoClient(ctx)
	dataCollection := client.Database("Lunchbox").Collection("UserData")
	result := dataCollection.FindOne(ctx, database.UserFavouriteRestaurants{UserId: user.ID})

	if result.Err() != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": result.Err().Error(),
		})
		return
	}
	var favs database.UserFavouriteRestaurants
	if err := result.Decode(&favs); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"result": favs.FavResIds,
	})
}
