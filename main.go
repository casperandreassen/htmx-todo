package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

var db *sqlx.DB

func main() {
	database, err := sqlx.Connect("sqlite3", "./todos.db")
	if err != nil {
		log.Fatalln(err)
	}
	db = database
	router := gin.Default()
	//load html file
	router.LoadHTMLGlob("templates/*.html")

	//static path
	router.Static("/assets", "./assets")
	router.GET("/", handleIndex)
	router.GET("/todo", handleGetTodos)
	router.POST("/todo", handleNewTodo)
	router.DELETE("/todo", handleDeleteTodo)
	router.PATCH("/todo", handleUpdateTodoState)
	router.Run(":8080")

}

type DBTodo struct {
	Id          int          `db:"id"`
	Status      int          `db:"status"`
	Title       string       `form:"title" binding:"required" db:"title"`
	Description string       `form:"description" binding:"required" db:"description"`
	Date        sql.NullTime `form:"date" db:"date"`
}

type Todo struct {
	Status      int    `db:"status"`
	Title       string `form:"title" binding:"required" db:"title"`
	Description string `form:"description" binding:"required" db:"description"`
	Date        string `form:"date" db:"date"`
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func handleGetTodos(c *gin.Context) {
	todos := []DBTodo{}
	db.Select(&todos, "SELECT * FROM todos")
	completedTodos := []DBTodo{}
	expiredTodos := []DBTodo{}
	otherTodos := []DBTodo{}
	for i := range todos {
		if todos[i].Status == 1 {
			completedTodos = append(completedTodos, todos[i])
		} else if todos[i].Date.Valid {

			if time.Now().After(todos[i].Date.Time) {
				expiredTodos = append(expiredTodos, todos[i])
			} else {
				otherTodos = append(otherTodos, todos[i])
			}
		} else {
			otherTodos = append(otherTodos, todos[i])
		}
	}
	c.HTML(http.StatusOK, "todos.html", gin.H{"completedTodos": completedTodos, "expiredTodos": expiredTodos, "otherTodos": otherTodos})
}

func handleNewTodo(c *gin.Context) {
	var data Todo
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	db.MustExec("INSERT INTO todos (title, description, status, date) VALUES ($1, $2, $3, $4)", data.Title, data.Description, 0, NewNullString(data.Date))
	todos := []DBTodo{}
	db.Select(&todos, "SELECT * FROM todos")

	handleGetTodos(c)
}

func handleDeleteTodo(c *gin.Context) {
	type Id struct {
		Id int `form:"id"`
	}
	data := Id{}
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	db.MustExec("DELETE FROM todos WHERE id = $1", data.Id)

	handleGetTodos(c)
}

func handleUpdateTodoState(c *gin.Context) {
	type Id struct {
		Id int `form:"id"`
	}
	data := Id{}
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	todo := DBTodo{}

	db.Get(&todo, "SELECT * FROM todos WHERE todos.id = $1", data.Id)

	var newStatus int = 0
	if todo.Status == 0 {
		newStatus = 1
	}

	db.MustExec("UPDATE todos SET status = $1 where id = $2", newStatus, todo.Id)

	handleGetTodos(c)
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
