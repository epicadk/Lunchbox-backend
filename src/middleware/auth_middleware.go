package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//TokenAuthMiddleware to verify AuthToken
func TokenAuthMiddleware(c *gin.Context) {

	token := c.Request.Header["Auth_token"][0]
	if token == "" {
		respondWithError(c, http.StatusUnauthorized, "Not Authorized")
		return
	}

	verified, err := verifyAuthToken(token)
	if err != nil {
		if err.Error() == "Token is expired" {
			respondWithError(c, http.StatusUnauthorized, err.Error())
		} else {
			respondWithError(c, http.StatusUnauthorized, "Invalid API Token")
		}
		return
	}

	if !verified {
		respondWithError(c, http.StatusUnauthorized, "Invalid API token")
		return
	}

	c.Next()

}

//verifyAuthToken Validates authtoken
func verifyAuthToken(header string) (bool, error) {

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

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": message,
	})
}
