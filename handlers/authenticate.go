package handlers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-server/db"
	"go-server/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
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
	user, err := getUser(credentials.Username)

	if err != nil {
		c.HTML(http.StatusOK, "login", gin.H{"errorMessage": "No such user."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.HTML(http.StatusOK, "login", gin.H{"errorMessage": "Invalid credentials."})
		return
	} else {
		token, err := utils.IssueJwtToken(user.Id, user.Username)
		if err == nil {
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("token", token, 6000000, "/", "htmx-todo-23.fly.dev", true, true)
			c.Header("HX-Redirect", "/")
		} else {
			c.HTML(http.StatusOK, "login", gin.H{"errorMessage": "Something went wrong."})
		}
		return
	}
}

func SignOut(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	app_URL := os.Getenv("APP_URL")
	c.SetCookie("token", "", 31556926, "/", app_URL, true, true)
	c.Header("HX-Redirect", "/login")
}

func getUser(username string) (db.User, error) {
	var user db.User

	row := db.DB.QueryRow("SELECT * FROM user WHERE username = ?", username)
	if err := row.Scan(&user.Id, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("No such row")
		}
		return user, errors.New("Error getting todo")
	}
	return user, nil
}
