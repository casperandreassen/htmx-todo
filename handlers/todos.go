package handlers

import (
	"github.com/gin-gonic/gin"
	"go-server/db"
	"go-server/utils"
	"net/http"
)

type Todo struct {
	Status      int    `db:"status"`
	Title       string `form:"title" binding:"required" db:"title"`
	Description string `form:"description" binding:"required" db:"description"`
	Date        string `form:"date" db:"date"`
}

func HandleGetTodos(c *gin.Context) {
	todos, err := db.GetTodos(c)
	if err != nil {
		c.AbortWithStatus(500)
	}
	completedTodos, expiredTodos, otherTodos, err := utils.TransformTodos(todos)
	if err != nil {
		c.AbortWithStatus(500)
	}
	c.HTML(http.StatusOK, "todos_page", gin.H{"completedTodos": completedTodos, "expiredTodos": expiredTodos, "otherTodos": otherTodos})
}

func HandleGetTodoElements(c *gin.Context) {
	todos, err := db.GetTodos(c)
	if err != nil {
		c.AbortWithStatus(500)
	}
	completedTodos, expiredTodos, otherTodos, err := utils.TransformTodos(todos)
	if err != nil {
		c.AbortWithStatus(500)
	}
	c.HTML(http.StatusOK, "todos", gin.H{"completedTodos": completedTodos, "expiredTodos": expiredTodos, "otherTodos": otherTodos})
}

func HandleNewTodo(c *gin.Context) {
	userid := c.MustGet("id").(float64)
	var data Todo
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	_, err := db.InsertTodo(data, int(userid))

	if err != nil {
		c.AbortWithStatus(500)
		return
	}

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
	db.DeleteTodo(data.Id)
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

	todo, err := db.GetSingleTodo(data.Id)

	if err != nil {
		c.AbortWithStatus(500)
	}

	var newStatus int = 0
	if todo.Status == 0 {
		newStatus = 1
	}

	err = db.UpdateTodoStatus(data.Id, newStatus)

	if err != nil {
		c.AbortWithStatus(500)
	}

	HandleGetTodoElements(c)
}
