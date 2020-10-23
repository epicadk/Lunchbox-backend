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
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Comment  string             `json:"comment,omitempty" bson:"comment,omitempty"`
}

//CommentsPost handles Post Request on comment Endpoint
func CommentsPost(c *gin.Context) {
	user := c.PostForm("UserName")
	comment := c.PostForm("Comment")
	resID := c.PostForm("ResID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	var curruser User
	var commentContainer CommentContainer
	commentContainer.Username = user
	commentContainer.Comment = comment
	usercollection := client.Database("Lunchbox").Collection("Users")
	err := usercollection.FindOne(ctx, User{Username: user}).Decode(&curruser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Username does not exist",
			"error":   err.Error(),
		})
		return
	}
	commentCollection := client.Database("LunchboxComments").Collection(resID)
	_, err = commentCollection.InsertOne(ctx, commentContainer)
	c.JSON(200, gin.H{
		"message": "done",
		"comment": commentContainer.Comment,
	})

}
