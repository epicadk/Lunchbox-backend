package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
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
	token := c.GetHeader("Auth_token")
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}
	userID := utils.GetUserID(token)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists

	userCollection := client.Database("Lunchbox").Collection("Users")
	result := userCollection.FindOne(ctx, database.User{ID: objectID})
	if result.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
			"error":   result.Err(),
		})
		return
	}

	//Decoding user

	var user database.User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
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
	res := favsCollection.FindOneAndUpdate(ctx, UserFavs{UserId: objectID}, update, &opt)

	if res.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist result error",
			"error":   res.Err().Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": http.StatusText(200),
	})

}
