package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//LoginRequest model
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//LoginResponse model
type LoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

//LoginPost Handles post Request to Login Endpoint
func LoginPost(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.Username == "" || request.Password == "" {
		utils.RespondWithQuickError(c, 400)
		return
	}
	//TODO add password validation here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := database.CheckUserFromUsername(ctx, c, request.Username)
	if err != nil {
		return
	}

	if !comparePasswords(user.Password, request.Password) {
		utils.RespondWithError(c, http.StatusUnauthorized, "Password is incorrect")
		return
	}

	//Generating Auth Token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
	}
	RefreshToken, err := utils.GenerateRefreshJWT(user.ID, user.Password)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}

	//Final Response if all is okay
	response := LoginResponse{AuthToken: token,
		RefreshToken: RefreshToken}
	c.JSON(200, response)
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}
