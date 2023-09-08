package handlers

import (
	"errors"
	"go-server/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	if credentials.Password != credentials.RePassword {
		c.HTML(201, "signup", gin.H{"errorMessage": "Passwords does not match"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)

	if err != nil {
		c.Status(500)
	}

	var stringHash = string(hash)

	_, err = insertUser(credentials, stringHash)

	if err != nil {
		c.HTML(201, "signup", gin.H{"errorMessage": "Username already taken"})
		return
	}
	c.HTML(201, "account_created.html", gin.H{})
}

func insertUser(credentials RegisterUserInput, passwordHash string) (int64, error) {
	result, err := db.DB.Exec("INSERT INTO user (username, password) VALUES (?, ?)", credentials.Username, passwordHash)
	if err != nil {
		return 0, errors.New("Could not insert user")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("Could not get id of inserted row")
	}
	return id, nil
}
