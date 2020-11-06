package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func RecentActivityGet(c *gin.Context) {
	username := c.Query("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := database.CheckUserFromUsername(ctx, c, username)

	if err != nil {
		return
	}

	client := database.MongoClient(ctx)
	commentsCollection := client.Database("Lunchbox").Collection("Comments")
	option := options.Find()
	option.SetLimit(10)
	cursor, err := commentsCollection.Find(ctx, database.CommentsContainer{UserID: user.ID}, option)
	if err != nil {
		utils.RespondWithQuickError(c, 204)
	}

	defer cursor.Close(ctx)
	var comments []database.CommentsContainer
	for cursor.Next(ctx) {
		var comment database.CommentsContainer
		if err = cursor.Decode(&comment); err != nil {
			utils.RespondWithQuickError(c, 500)
			return
		}
		comments = append(comments, comment)
	}

	c.JSON(200, gin.H{
		"comments": comments,
	})

}
