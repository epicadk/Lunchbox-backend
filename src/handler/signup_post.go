package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

//User defines user
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

//TODO add Validation

//SignupPost handles post request to Signup End Point
func SignupPost(c *gin.Context) {

	//Getting data form the Request
	username := c.PostForm("username")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")

	//Valid password and repassword
	if password != repassword {
		c.JSON(400, gin.H{
			"message": "passwords do not match",
		})
		return
	}

	// Connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		//Add mongodb URI here
		"",
	))
	if err != nil {
		log.Fatal(err)
	}

	// Checking if user Exists
	var user User
	collection := client.Database("Lunchbox").Collection("Users")
	err = collection.FindOne(ctx, User{Username: username}).Decode(&user)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Username exists",
		})
		return
	}

	// If user does not exists add user to database and return info
	user.Username = username
	user.Password = hashpassword(password)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	//Generate Auth Token For current User
	token, err := GenerateJWT(result.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User was created successfully.",
		"id":      token,
	})
}

func hashpassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
