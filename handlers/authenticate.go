package handlers

import (
	"github.com/gin-gonic/gin"
	"go-server/db"
	"go-server/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateInput struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Authenticate(c *gin.Context) {
	var credentials AuthenticateInput

	if err := c.ShouldBind(&credentials); err != nil {
		return
	}
	retrievedUser := db.User{}

	db.DB.Get(&retrievedUser, "SELECT * FROM user WHERE username = $1", credentials.Username)

	err := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(credentials.Password))
	if err != nil {
		c.HTML(401, "invalid_credentials.html", gin.H{"errorMessage": "Invalid credentials"})
		return
	} else {
		token, err := utils.IssueJwtToken(retrievedUser)
		if err == nil {
			c.SetCookie("token", token, 6000000, "/", "htmx-todo.fly.dev", true, true)
			c.Header("HX-Redirect", "/")

		} else {
			c.Status(500)
		}
		return
	}
}
