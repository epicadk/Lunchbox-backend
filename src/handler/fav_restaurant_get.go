package handler

import (
	"context"
	"go-gin-api/src/database"
	"log"
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
	client := database.MongoClient(ctx)
	userCollection := client.Database("Lunchbox").Collection("Users")
	result := userCollection.FindOne(ctx, database.User{Username: username})

	if result.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
			"error":   result.Err(),
			"result":  result,
		})
		return
	}
	//Decoding user
	var user database.User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}

	dataCollection := client.Database("Lunchbox").Collection("UserData")
	result = dataCollection.FindOne(ctx, database.UserFavouriteRestaurants{UserId: user.ID})

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
