package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserFavouritePostRequest struct {
	Favourite int64 `json:"zomato_res_id"`
}

//TODO Delete this
type UserFavs struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FavResIds []int64            `json:"fav_res_ids,omitempty" bson:"fav_res_ids,omitempty"`
	UserId    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

//FavRestaurantPost handles post req to post endpoint.
func FavRestaurantPost(c *gin.Context) {

	var request UserFavouritePostRequest
	token := c.GetHeader("Authorization")
	if err := c.ShouldBindJSON(&request); err != nil || request.Favourite == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}
	userID, err := utils.GetUserID(token)
	//Possible only when Someone has cracked the secret
	if err != nil {
		utils.RespondWithError(c, 401, "Invalid AuthToken")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	_, err = database.CheckUserFromUserID(ctx, c, userID)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}

	//InsertComments in database

	favsCollection := client.Database("Lunchbox").Collection("UserData")
	update := bson.M{
		"$addToSet": bson.M{"fav_res_ids": request.Favourite},
	}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}
	res := favsCollection.FindOneAndUpdate(ctx, UserFavs{UserId: userID}, update, &opt)

	if res.Err() != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}
	c.JSON(200, gin.H{
		"message": http.StatusText(200),
	})

}
