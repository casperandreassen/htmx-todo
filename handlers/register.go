package handlers

import (
	"github.com/gin-gonic/gin"
	"go-server/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RegisterUserInput struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	RePassword string `form:"re-password"`
}

func RegisterUser(c *gin.Context) {
	var credentials RegisterUserInput

	if err := c.ShouldBind(&credentials); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)

	if err != nil {
		c.Status(500)
	}

	var stringHash = string(hash)

	db.DB.MustExec("INSERT INTO user (username, password) VALUES ($1, $2)", "casper", stringHash)

	c.HTML(201, "account_created.html", gin.H{})

}
