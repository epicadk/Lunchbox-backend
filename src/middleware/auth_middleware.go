package middleware

import (
	"github.com/epicadk/Lunchbox-backend/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//TokenAuthMiddleware to verify AuthToken
func TokenAuthMiddleware(c *gin.Context) {

	//TODO check if Auth header does not exist.
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		utils.RespondWithQuickError(c, http.StatusUnauthorized)
		return
	}

	verified, err := utils.VerifyAuthToken(token)
	if err != nil {
		if err.Error() == "Token is expired" {
			utils.RespondWithError(c, http.StatusUnauthorized, err.Error())
		} else {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid API Token")
		}
		return
	}

	if !verified {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid API token")
		return
	}

	c.Next()

}
