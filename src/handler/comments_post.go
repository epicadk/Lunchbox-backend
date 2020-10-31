package handler

import (
	"context"
	database "go-gin-api/src/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CommentContainer comment Model
type CommentContainer struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Comment     string             `json:"comment" bson:"comment"`
	ZomatoResID string             `json:"zomato_res_id" bson:"zomato_res_id"`
}

//CommentsPost handles Post Request on comment Endpoint
func CommentsPost(c *gin.Context) {

	userID := c.PostForm("UserID")
	comment := c.PostForm("Comment")
	ZomatoResID := c.PostForm("ZomatoResID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	var commentContainer CommentContainer
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": http.StatusText(http.StatusUnauthorized),
		})
	}

	commentContainer.UserID = objectID
	commentContainer.Comment = comment
	commentContainer.ZomatoResID = ZomatoResID

	// Checking if user Exists
	userCollection := client.Database("Lunchbox").Collection("Users")
	result := userCollection.FindOne(ctx, User{ID: commentContainer.UserID})

	if result.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
		})
		return
	}

	//Decoding user
	var user User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}

	//Adding comment to db
	commentCollection := client.Database("Lunchbox").Collection("Comments")
	_, err = commentCollection.InsertOne(ctx, commentContainer)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": http.StatusText(500),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "done",
		"comment": commentContainer.Comment,
	})

}
