package handler

import (
	"context"
	"fmt"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type CommentsGetResponse struct {
	Comments    []database.CommentsContainer `json:"comments"`
	UserComment database.CommentsContainer   `json:"user_comment,omitempty"`
}

func CommentsGet(c *gin.Context) {
	tresId := c.Query("resID")
	resId, err := strconv.Atoi(tresId)
	if err != nil {
		utils.RespondWithQuickError(c, 400)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	userComment, err := database.FindCommentByUserID(ctx, c, c.GetHeader("Authorization"), resId)
	if err != nil {
		return
	}
	commentsCollection := client.Database("Lunchbox").Collection("Comments")
	cursor, err := commentsCollection.Find(ctx, database.CommentsContainer{ZomatoResID: resId})

	if err != nil {
		utils.RespondWithQuickError(c, 204)
		return
	}

	defer cursor.Close(ctx)
	var comments []database.CommentsContainer
	for cursor.Next(ctx) {
		var comment database.CommentsContainer
		if err = cursor.Decode(&comment); err != nil {
			utils.RespondWithQuickError(c, 500)
			return
		}
		fmt.Println(comments)
		comments = append(comments, comment)
	}

	c.JSON(200, CommentsGetResponse{Comments: comments, UserComment: userComment})
}
