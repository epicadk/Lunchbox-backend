package database

import (
	"context"
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

//ConnectDB lol
func MongoClient(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		//Add mongodb URI here
		"",
	))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func CheckUserFromUsername(ctx context.Context, c *gin.Context, Username string) (User, error) {

	client := MongoClient(ctx)
	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, User{Username: Username})
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
		return User{}, result.Err()
	}

	var user User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}
	return user, nil
}

func CheckUserFromUserID(ctx context.Context, c *gin.Context, UserID primitive.ObjectID) (User, error) {

	client := MongoClient(ctx)
	collection := client.Database("Lunchbox").Collection("Users")
	result := collection.FindOne(ctx, User{ID: UserID})
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
		return User{}, result.Err()
	}

	var user User
	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
	}
	return user, nil
}

func GetCommentByCommentId(ctx context.Context, c *gin.Context, commentID primitive.ObjectID) (CommentsContainer, error) {
	client := MongoClient(ctx)
	var result CommentsContainer

	commentsCollection := client.Database("Lunchbox").Collection("Comments")
	res := commentsCollection.FindOne(ctx, CommentsContainer{ID: commentID})
	if res.Err() != nil {
		utils.RespondWithQuickError(c, http.StatusNoContent)
		return result, res.Err()
	}
	if err := res.Decode(result); err != nil {
		utils.RespondWithQuickError(c, http.StatusInternalServerError)
		return result, err
	}
	return result, nil
}

func DeleteCommentByCommentId(ctx context.Context, c *gin.Context, commentID primitive.ObjectID, token string) bool {
	client := MongoClient(ctx)
	userid, err := utils.GetUserID(token)

	if err != nil {
		utils.RespondWithQuickError(c, http.StatusInternalServerError)
		return false
	}

	commentsCollection := client.Database("Lunchbox").Collection("Comments")
	res, err := commentsCollection.DeleteOne(ctx, CommentsContainer{ID: commentID, UserID: userid})

	if err != nil {
		utils.RespondWithQuickError(c, http.StatusInternalServerError)
		return false
	}
	if res.DeletedCount != 1 {
		utils.RespondWithQuickError(c, http.StatusNoContent)
		return false
	}
	return true
}

func UpdateCommentByCommentId(ctx context.Context, c *gin.Context, commentID primitive.ObjectID, token string, update bson.M) bool {
	client := MongoClient(ctx)
	userid, err := utils.GetUserID(token)

	if err != nil {
		utils.RespondWithQuickError(c, http.StatusInternalServerError)
		return false
	}

	commentsCollection := client.Database("Lunchbox").Collection("Comments")
	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}
	res := commentsCollection.FindOneAndUpdate(ctx, CommentsContainer{ID: commentID, UserID: userid}, update, &opt)
	if res.Err() != nil {
		utils.RespondWithQuickError(c, http.StatusInternalServerError)
		return false

	}
	return true
}
