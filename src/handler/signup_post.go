package handler

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"//Add Mongo URI here",
	))
	if err != nil {
		log.Fatal(err)
	}

	// Checking if user Exists
	var user User
	collection := client.Database("Lunchbox").Collection("Users")
	err = collection.FindOne(ctx, User{Username: username}).Decode(&user)
	if err == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Username exists",
		})
		return
	}

	// If user does not exists add user to database and return info
	user.Username = username
	user.Password = password
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	token, err := GenerateJWT(result.InsertedID.(primitive.ObjectID).String())
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message":    user.Username,
		"password":   user.Password,
		"repassword": repassword,
		"id":         token,
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
	tokenString, err := token.SignedString([]byte("//Add Secret Here"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
