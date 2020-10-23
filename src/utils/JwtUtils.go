package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

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
	return false, nil
}
