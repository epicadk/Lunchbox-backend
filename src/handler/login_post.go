package handler

import (
	"context"
	database "go-gin-api/src/database"
	"go-gin-api/src/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//Login request model
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login response model
type LoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

//LoginPost Handles post Request to Login Endpoint
func LoginPost(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": http.StatusText(http.StatusUnprocessableEntity),
			"error":   err.Error(),
		})
		return
	}
	//TODO add password validation here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, database.User{Username: request.Username})

	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "Username does not exist",
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

	if !comparePasswords(user.Password, request.Password) {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Password incorrect",
		})
		return
	}
	//Generating Auth Token
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": http.StatusText(500),
		})
		log.Fatal(err.Error())
	}

	//Final Response if all is okay
	response := LoginResponse{AuthToken: token,
		RefreshToken: ""}
	c.JSON(200, response)
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}
