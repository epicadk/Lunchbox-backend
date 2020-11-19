package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type CommentDeleteRequest struct {
	CommentId string `json:"commentid"`
}

func CommentsDelete(c *gin.Context) {
	var request CommentDeleteRequest
	token := c.GetHeader("Authorization")
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondWithQuickError(c, http.StatusBadRequest)
		return
	}
	commentid, err := primitive.ObjectIDFromHex(request.CommentId)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid Comment ID")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	success := database.DeleteCommentByCommentId(ctx, c, commentid, token)
	if !success {
		return
	}
	c.JSON(204, gin.H{
		"message": "Comment Deleted",
	})
}
