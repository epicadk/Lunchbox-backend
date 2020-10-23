package handler

import (
	"context"
	database "go-gin-api/src/database"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//LoginPost Handles post Request to Login Endpoint
func LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	var user User
	collection := client.Database("Lunchbox").Collection("Users")
	err := collection.FindOne(ctx, User{Username: username}).Decode(&user)

	if err != nil || user.Username == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username does not exist",
			"error":   err.Error(),
		})
		return
	}

	if !comparePasswords(user.Password, password) {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Password incorrect",
		})
		return
	}
	token, err := GenerateJWT(user.ID.Hex())

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Internal Server Error",
		})
		log.Fatal(err.Error())
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

//GenerateJWT generates JWT token
func GenerateJWT(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	//TODO add secret here
	tokenString, err := token.SignedString([]byte("mysecret"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
