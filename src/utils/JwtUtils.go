package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//GenerateJWT generates JWT token
func GenerateJWT(id primitive.ObjectID) (string, error) {
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

func GenerateRefreshJWT(id primitive.ObjectID, hashedPass string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = id
	claims["exp"] = time.Now().Add(time.Hour * 336).Unix()
	//TODO add secret here
	tokenString, err := token.SignedString([]byte("mysecret" + hashedPass))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyRefreshToken(header, hashedPass string) (bool, error) {
	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		// TODO
		return []byte("mysecret" + hashedPass), nil
	})
	if err != nil {
		return false, err
	}
	if token.Valid {
		return true, nil
	}
	log.Println(token)
	return false, nil
}

//VerifyAuthToken Validates authtoken
func VerifyAuthToken(header string) (bool, error) {

	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		// TODO
		return []byte("mysecret"), nil
	})
	if err != nil {
		return false, err
	}
	if token.Valid {
		return true, nil
	}
	log.Println(token)
	return false, nil
}

func GetUserID(header string) (primitive.ObjectID, error) {
	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(header, claims)
	if err != nil {
		return [12]byte{}, err
	}
	return primitive.ObjectIDFromHex(fmt.Sprint(claims["client"]))
}
