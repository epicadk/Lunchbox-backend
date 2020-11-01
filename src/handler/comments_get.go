package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin-api/src/database"
	"log"
	"net/http"
	"time"
)

type CommentsGetRequests struct {
	ZomatoResID string `json:"zomato_res_id"`
}

func CommentsGet(c *gin.Context) {
	var request CommentsGetRequests

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	commentsCollection := client.Database("Lunchbox").Collection("Comments")
	cursor, err := commentsCollection.Find(ctx, database.CommentsContainer{ZomatoResID: request.ZomatoResID})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)
	var comments []database.CommentsContainer
	for cursor.Next(ctx) {
		var comment database.CommentsContainer
		if err = cursor.Decode(&comment); err != nil {
			log.Fatal(err)
		}
		fmt.Println(comments)
		comments = append(comments, comment)
	}

	c.JSON(200, gin.H{
		"comments": comments,
	})
}
