package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go-server/db"
	"go-server/utils"
	"net/http"
	"time"
)

type DBTodo struct {
	Id          int          `db:"id"`
	Status      int          `db:"status"`
	Title       string       `form:"title" binding:"required" db:"title"`
	Description string       `form:"description" binding:"required" db:"description"`
	Date        sql.NullTime `form:"date" db:"date"`
	Userid      int          `db:"userid"`
}

type Todo struct {
	Status      int    `db:"status"`
	Title       string `form:"title" binding:"required" db:"title"`
	Description string `form:"description" binding:"required" db:"description"`
	Date        string `form:"date" db:"date"`
}

func HandleGetTodos(c *gin.Context) {
	completedTodos, expiredTodos, otherTodos := getTodos(c)
	c.HTML(http.StatusOK, "todos_page", gin.H{"completedTodos": completedTodos, "expiredTodos": expiredTodos, "otherTodos": otherTodos})
}

func HandleGetTodoElements(c *gin.Context) {
	completedTodos, expiredTodos, otherTodos := getTodos(c)
	c.HTML(http.StatusOK, "todos", gin.H{"completedTodos": completedTodos, "expiredTodos": expiredTodos, "otherTodos": otherTodos})
}

func HandleNewTodo(c *gin.Context) {
	userid := c.MustGet("id").(float64)
	var data Todo
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	db.DB.MustExec("INSERT INTO todos (title, description, status, date, userid) VALUES ($1, $2, $3, $4, $5)", data.Title, data.Description, 0, utils.NewNullString(data.Date), userid)
	todos := []DBTodo{}
	db.DB.Select(&todos, "SELECT * FROM todos")

	HandleGetTodoElements(c)
}

func HandleDeleteTodo(c *gin.Context) {
	type Id struct {
		Id int `form:"id"`
	}
	data := Id{}
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	db.DB.MustExec("DELETE FROM todos WHERE id = $1", data.Id)

	HandleGetTodoElements(c)
}

func HandleUpdateTodoState(c *gin.Context) {
	type Id struct {
		Id int `form:"id"`
	}
	data := Id{}
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	todo := DBTodo{}

	db.DB.Get(&todo, "SELECT * FROM todos WHERE todos.id = $1", data.Id)

	var newStatus int = 0
	if todo.Status == 0 {
		newStatus = 1
	}

	db.DB.MustExec("UPDATE todos SET status = $1 where id = $2", newStatus, todo.Id)

	HandleGetTodoElements(c)
}

func getTodos(c *gin.Context) ([]DBTodo, []DBTodo, []DBTodo) {
	userid := c.MustGet("id").(float64)
	todos := []DBTodo{}
	db.DB.Select(&todos, "SELECT * FROM todos WHERE userid = $1", int(userid))
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
	return completedTodos, expiredTodos, otherTodos
}
