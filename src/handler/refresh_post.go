package handler

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
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
	id, err := utils.GetUserID(RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Invalid Refresh Token",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := database.CheckUserFromUserID(ctx, c, id)
	if err != nil {
		return
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

	newToken, err := utils.GenerateJWT(user.ID)

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
