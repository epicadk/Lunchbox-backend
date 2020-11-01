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
	RefreshToken string `json:"refresh_token"`
}

//LoginPost Handles post Request to Login Endpoin
func LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	//TODO add should Bind JSON
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, User{Username: username})

	if result.Err() != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
			"error":   result.Err(),
			"result":  result,
		})
		return
	}
	//Decoding user
	var user User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}

	if !comparePasswords(user.Password, password) {
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
	c.JSON(200, gin.H{
		"token": token,
	})
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
