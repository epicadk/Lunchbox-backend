package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/epicadk/Lunchbox-backend/src/database"
	"github.com/epicadk/Lunchbox-backend/src/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//Request model is same as database user model

type SignupResponse struct {
	Message      string `json:"message"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

//TODO add Validation

//SignupPost handles post request to Signup End Point
func SignupPost(c *gin.Context) {
	var user database.User
	err := c.ShouldBindJSON(&user)
	if err != nil || user.Username == "" || user.Phone == 0 || user.Password == "" {
		utils.RespondWithQuickError(c, 400)
		return
	}
	//TODO add password validation
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		utils.RespondWithError(c, 400, err.Error())
		return
	}
	// Connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := database.MongoClient(ctx)

	// Checking if user Exists
	collection := client.Database("Lunchbox").Collection("Users")

	if result := collection.FindOne(ctx, database.User{Username: user.Username}); result.Err() == nil {
		utils.RespondWithError(c, http.StatusConflict, "Username already exists")
		return
	}
	// If user does not exists add user to database and return info
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	//Generate Auth Token For current User
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}
	RefreshToken, err := utils.GenerateRefreshJWT(user.ID, user.Password)
	if err != nil {
		utils.RespondWithQuickError(c, 500)
		return
	}
	c.JSON(http.StatusCreated, SignupResponse{Message: "User Created", AuthToken: token, RefreshToken: RefreshToken})
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
