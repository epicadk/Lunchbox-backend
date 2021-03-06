package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
	Phone    int64              `json:"phone" bson:"phone,omitempty"`
}

type UserFavouriteRestaurants struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FavResIds []int              `json:"fav_res_ids,omitempty" bson:"fav_res_ids,omitempty"`
	UserId    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type CommentsContainer struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Comment     string             `json:"comment" bson:"comment,omitempty"`
	Rating      int                `json:"rating" bson:"rating,omitempty"`
	ZomatoResID int                `json:"zomato_res_id" bson:"zomato_res_id,omitempty"`
	UserName    string             `json:"user_name,omitempty" bson:"user_name,omitempty"`
}
