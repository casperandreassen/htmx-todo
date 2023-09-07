package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go-server/utils"
)

func Auth(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.Abort()
		return
	}

	parsedToken, tokenErr := utils.VerifyJwtToken(cookie)

	if tokenErr != nil {
		c.Abort()
		return
	}

	if !parsedToken.Valid {
		c.Abort()
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})
	id := data["id"].(float64)
	level := data["username"].(string)
	if ok && parsedToken.Valid {
		c.Set("id", id)
		c.Set("username", level)
		c.Next()
	} else {
		c.Abort()
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

	if !parsedToken.Valid {
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})
	id := data["id"].(float64)
	level := data["username"].(string)
	if ok && parsedToken.Valid {
		c.Set("id", id)
		c.Set("username", level)
		return true
	} else {
		return false
	}
}