package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type CommentPutRequest struct {
	CommentID   string `json:"comment_id"`
	CommentBody string `json:"comment_body"`
	Rating      int    `json:"rating"`
}

func CommentsPut(c *gin.Context) {
	var request CommentPutRequest
	token := c.GetHeader("Authorization")

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondWithQuickError(c, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	commentId, err := primitive.ObjectIDFromHex(request.CommentID)
	if err != nil {
		utils.RespondWithQuickError(c, http.StatusInternalServerError)
		return
	}

	update := bson.M{
		"$set": bson.M{"rating": request.Rating, "comment": request.CommentBody},
	}

	success := database.UpdateCommentByCommentId(ctx, c, commentId, token, update)
	if !success {
		return
	}
	c.JSON(200, gin.H{
		"message": "Comment Updated",
	})

}
