package database

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

//ConnectDB lol
func MongoClient(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		//Add mongodb URI here
		"",
	))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func CheckUserFromUsername(ctx context.Context, c *gin.Context, Username string) (User, error) {

	client := MongoClient(ctx)
	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, User{Username: Username})
	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "Username does not exist",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return User{}, result.Err()
	}

	var user User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}
	return user, nil
}

func CheckUserFromUserID(ctx context.Context, c *gin.Context, UserID primitive.ObjectID) (User, error) {

	client := MongoClient(ctx)
	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, User{ID: UserID})
	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "Username does not exist",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return User{}, result.Err()
	}

	var user User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}
	return user, nil
}
