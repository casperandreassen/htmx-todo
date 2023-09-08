package middleware

import (
	"github.com/gin-gonic/gin"
	"go-server/utils"
	"net/http"
)

func Auth(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	parsedToken, tokenErr := utils.VerifyJwtToken(cookie)

	if tokenErr != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !parsedToken.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := parsedToken.Claims.(*utils.TodoClaims); ok && parsedToken.Valid {
		c.Set("id", claims.Userid)
		c.Set("username", claims.Username)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func IsUserSignedIn(c *gin.Context) bool {
	cookie, err := c.Cookie("token")
	if err != nil {
		return false
	}
	parsedToken, tokenErr := utils.VerifyJwtToken(cookie)

	if tokenErr != nil {
		return false
	}
	return parsedToken.Valid
}
