package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	database "go-gin-api/src/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type UserFavs struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FavResIds []int              `json:"fav_res_ids" bson:"fav_res_ids"`
	//Check this
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
}

//FavRestaurantPost handles post req to post endpoint.
func FavRestaurantPost(c *gin.Context) {
	var userFavs UserFavs
	var favid int
	err := c.ShouldBindJSON(&userFavs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	var user User
	userCollection := client.Database("Lunchbox").Collection("Users")
	err = userCollection.FindOne(ctx, User{ID: userFavs.UserId}).Decode(&user)

	if err != nil || user.Username == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
			"error":   err.Error(),
		})
		return
	}

	//InsertComments in database

	favsCollection := client.Database("Lunchbox").Collection("UserData")
	update := bson.M{
		"$addToSet": bson.M{"fav_res_ids": favid},
	}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}
	result := favsCollection.FindOneAndUpdate(ctx, UserFavs{UserId: userFavs.UserId}, update, &opt)

	if result.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
			"error":   result.Err(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": http.StatusText(200),
	})

}
