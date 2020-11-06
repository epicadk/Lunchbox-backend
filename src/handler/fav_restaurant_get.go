package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func FavouriteRestaurantGet(c *gin.Context) {
	//TODO change to url param
	username := c.Query("username")
	if username == "" {
		utils.RespondWithQuickError(c, 400)
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
		utils.RespondWithError(c, http.StatusNoContent, "No Favourite Restaurants Found For User")
		return
	}
	var favs database.UserFavouriteRestaurants
	if err := result.Decode(&favs); err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}

	c.JSON(200, gin.H{
		"result": favs.FavResIds,
	})
}
