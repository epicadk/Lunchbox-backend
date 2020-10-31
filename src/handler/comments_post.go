package handler

import (
	"context"
	database "go-gin-api/src/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CommentContainer comment Model
type CommentContainer struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username    string             `json:"username,omitempty" bson:"username,omitempty"`
	Comment     string             `json:"comment,omitempty" bson:"comment,omitempty"`
	ZomatoResID string             `json:"zomato_res_id,omitempty" bson:"zomato_res_id,omitempty"`
}

//CommentsPost handles Post Request on comment Endpoint
func CommentsPost(c *gin.Context) {

	username := c.PostForm("UserName")
	comment := c.PostForm("Comment")
	ZomatoResID := c.PostForm("ZomatoResID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	var user User
	var commentContainer CommentContainer

	commentContainer.Username = username
	commentContainer.Comment = comment
	commentContainer.ZomatoResID = ZomatoResID

	// Checking if user Exists
	userCollection := client.Database("Lunchbox").Collection("Users")
	err := userCollection.FindOne(ctx, User{Username: username}).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Username does not exist",
			"error":   err.Error(),
		})
		return
	}

	//Adding comment to db
	commentCollection := client.Database("Lunchbox").Collection("Comments")
	_, err = commentCollection.InsertOne(ctx, commentContainer)
	c.JSON(200, gin.H{
		"message": "done",
		"comment": commentContainer.Comment,
	})

}
