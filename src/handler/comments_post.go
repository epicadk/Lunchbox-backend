package handler

import (
	"context"
	database "go-gin-api/src/database"
	"go-gin-api/src/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentContainerPostRequest struct {
	Comment     string `json:"comment"`
	Title       string `json:"title"`
	ZomatoResID string `json:"zomato_res_id"`
}

//CommentContainer comment Model

//CommentsPost handles Post Request on comment Endpoint
func CommentsPost(c *gin.Context) {
	var request CommentContainerPostRequest
	token := c.GetHeader("Auth_token")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
			"error":   err.Error(),
		})
		return
	}

	userID := utils.GetUserID(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	var commentContainer database.CommentsContainer
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	commentContainer.UserID = objectID
	commentContainer.Comment = request.Comment
	commentContainer.ZomatoResID = request.ZomatoResID
	commentContainer.Title = request.Title

	// Checking if user Exists
	userCollection := client.Database("Lunchbox").Collection("Users")
	result := userCollection.FindOne(ctx, database.User{ID: commentContainer.UserID})

	if result.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
		})
		return
	}

	//Decoding user
	var user database.User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}

	commentContainer.UserName = user.Username

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
