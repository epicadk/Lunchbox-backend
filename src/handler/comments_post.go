package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"

	"github.com/gin-gonic/gin"
)

type CommentContainerPostRequest struct {
	Comment     string `json:"comment"`
	Rating      int    `json:"rating"`
	ZomatoResID int    `json:"zomato_res_id"`
}

//CommentContainer comment Model

//CommentsPost handles Post Request on comment Endpoint
func CommentsPost(c *gin.Context) {
	var request CommentContainerPostRequest
	token := c.GetHeader("Authorization")

	if err := c.ShouldBindJSON(&request); err != nil || request.Comment == "" || request.Rating == 0 || request.ZomatoResID == 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		utils.RespondWithError(c, 400, "Invalid AuthToken")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	var commentContainer database.CommentsContainer

	commentContainer.UserID = userID
	commentContainer.Comment = request.Comment
	commentContainer.ZomatoResID = request.ZomatoResID
	commentContainer.Rating = request.Rating

	user, err := database.CheckUserFromUserID(ctx, c, userID)
	if err != nil {
		return
	}

	commentContainer.UserName = user.Username

	//Adding comment to db
	commentCollection := client.Database("Lunchbox").Collection("Comments")
	_, err = commentCollection.InsertOne(ctx, commentContainer)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}
	c.JSON(200, gin.H{
		"message": "done",
		"comment": commentContainer.Comment,
	})

}
