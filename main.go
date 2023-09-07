package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go-server/db"
	"go-server/handlers"
	"go-server/middleware"
	"go-server/utils"
	"net/http"
)

var DB *sqlx.DB

func main() {
	utils.GenerateKeyPairs()
	db.InitDB()
	router := gin.Default()

	//static path
	router.Static("/assets", "./assets")
	router.GET("/login", HandleLoginPage)
	router.GET("/signup", HandleSignUpPage)
	router.POST("/authenticate", handlers.Authenticate)
	router.POST("/signup", handlers.RegisterUser)
	router.GET("/", HandleIndex)

	//load html file
	router.LoadHTMLGlob("templates/*.html")

	private := router.Group("/")
	{
		private.Use(middleware.Auth)
		private.GET("/todos", handlers.HandleGetTodos)
		private.GET("/todo", handlers.HandleGetTodoElements)
		private.POST("/todo", handlers.HandleNewTodo)
		private.DELETE("/todo", handlers.HandleDeleteTodo)
		private.PATCH("/todo", handlers.HandleUpdateTodoState)
	}
	router.Run(":8080")

}

func HandleIndex(c *gin.Context) {
	signedIn := middleware.IsUserSignedIn(c)
	c.HTML(http.StatusOK, "index.html", gin.H{"signedIn": signedIn})
}

func HandleSignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func HandleLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
