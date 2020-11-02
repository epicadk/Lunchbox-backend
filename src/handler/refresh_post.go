package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type RefreshResponse struct {
	AuthToken string `json:"auth_token"`
}

func RefreshPost(c *gin.Context) {
	RefreshToken := c.Request.Header.Get("Authorization")
	log.Println(c.Request.Header)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)
	id := utils.GetUserID(RefreshToken)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Refresh Token",
		})
		return
	}

	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, database.User{ID: objectId})

	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Invalid Refresh Token",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	//Decoding user
	var user database.User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}

	verified, err := utils.VerifyRefreshToken(RefreshToken, user.Password)
	if err != nil {
		if err.Error() == "Refresh Token is expired" {
			respondWithError(c, http.StatusUnauthorized, err.Error())
		} else {
			respondWithError(c, http.StatusUnauthorized, "Invalid Refresh Token")
		}
		return
	}
	if !verified {
		respondWithError(c, http.StatusUnauthorized, "Invalid Refresh token")
		return
	}

	newToken, err := utils.GenerateJWT(user.ID.String())

	if err != nil {
		respondWithError(c, 500, http.StatusText(500))
		return
	}
	c.JSON(200, RefreshResponse{AuthToken: newToken})
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"message": message,
	})
}
