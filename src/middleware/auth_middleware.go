package middleware

import (
	"go-gin-api/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//TokenAuthMiddleware to verify AuthToken
func TokenAuthMiddleware(c *gin.Context) {

	//TODO check if Auth header does not exist.
	token := c.Request.Header.Get("Auth_token")
	if token == "" {
		respondWithError(c, http.StatusUnauthorized, "Not Authorized")
		return
	}

	verified, err := utils.VerifyAuthToken(token)
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

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": message,
	})
}
