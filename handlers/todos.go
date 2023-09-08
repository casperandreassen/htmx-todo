package handlers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-server/db"
	"go-server/utils"
	"log"
	"net/http"
	"time"
)

type DBTodo struct {
	Id          int            `db:"id"`
	Status      int            `db:"status"`
	Title       string         `form:"title" binding:"required" db:"title"`
	Description string         `form:"description" binding:"required" db:"description"`
	Date        sql.NullString `form:"date" db:"date"`
	Userid      int            `db:"userid"`
}

type Todo struct {
	Status      int    `db:"status"`
	Title       string `form:"title" binding:"required" db:"title"`
	Description string `form:"description" binding:"required" db:"description"`
	Date        string `form:"date" db:"date"`
}

func HandleGetTodos(c *gin.Context) {
	todos, err := getTodos(c)
	if err != nil {
		c.AbortWithStatus(500)
	}
	completedTodos, expiredTodos, otherTodos, err := transformTodos(todos)
	if err != nil {
		c.AbortWithStatus(500)
	}
	c.HTML(http.StatusOK, "todos_page", gin.H{"completedTodos": completedTodos, "expiredTodos": expiredTodos, "otherTodos": otherTodos})
}

func HandleGetTodoElements(c *gin.Context) {
	todos, err := getTodos(c)
	if err != nil {
		c.AbortWithStatus(500)
	}
	completedTodos, expiredTodos, otherTodos, err := transformTodos(todos)
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

	_, err := insertTodo(data, int(userid))

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
	deleteTodo(data.Id)
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

	todo, err := getSingleTodo(data.Id)

	if err != nil {
		c.AbortWithStatus(500)
	}

	var newStatus int = 0
	if todo.Status == 0 {
		newStatus = 1
	}

	err = updateTodoStatus(data.Id, newStatus)

	if err != nil {
		c.AbortWithStatus(500)
	}

	HandleGetTodoElements(c)
}

func getTodos(c *gin.Context) ([]DBTodo, error) {
	userid := c.MustGet("id").(float64)
	todos := []DBTodo{}
	query, err := db.DB.Prepare("SELECT * FROM todos WHERE userid = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()
	rows, err := query.Query(int(userid))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var todo = DBTodo{}
		err := rows.Scan(&todo.Id, &todo.Status, &todo.Title, &todo.Description, &todo.Date, &todo.Userid)
		if err != nil {
			return nil, errors.New("Error scanning.")
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return todos, nil
}

func transformTodos(todos []DBTodo) ([]DBTodo, []DBTodo, []DBTodo, error) {
	completedTodos := []DBTodo{}
	expiredTodos := []DBTodo{}
	otherTodos := []DBTodo{}
	for i := range todos {
		if !todos[i].Date.Valid {
			if todos[i].Status == 1 {
				completedTodos = append(completedTodos, todos[i])
			} else {
				otherTodos = append(otherTodos, todos[i])
			}
		} else {
			todoDate, err := time.Parse("2006-01-02", todos[i].Date.String)
			if err != nil {
				return nil, nil, nil, errors.New("Could not parse date")
			}
			if todos[i].Status == 1 {
				completedTodos = append(completedTodos, todos[i])
			} else {
				if time.Now().After(todoDate) {
					expiredTodos = append(expiredTodos, todos[i])
				} else {
					otherTodos = append(otherTodos, todos[i])
				}
			}
		}

	}
	return completedTodos, expiredTodos, otherTodos, nil
}

func insertTodo(todo Todo, userid int) (int64, error) {
	result, err := db.DB.Exec("INSERT INTO todos (title, description, status, date, userid) VALUES (?, ?, ?, ?, ?)", todo.Title, todo.Description, 0, utils.NewNullString(todo.Date), userid)
	if err != nil {
		return 0, errors.New("Could not insert todo")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("Could not get id of inserted row")
	}
	return id, nil
}

func deleteTodo(todoId int) error {
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", todoId)
	if err != nil {
		return errors.New("Could not delete todo")
	}
	return nil
}

func updateTodoStatus(todoid int, newStatus int) error {
	_, err := db.DB.Exec("UPDATE todos SET status = ? where id = ?", newStatus, todoid)
	if err != nil {
		return errors.New("Could not update todo")
	}
	return nil
}

func getSingleTodo(todoid int) (DBTodo, error) {
	var todo DBTodo

	row := db.DB.QueryRow("SELECT * FROM todos WHERE todos.id = ?", todoid)
	if err := row.Scan(&todo.Id, &todo.Status, &todo.Title, &todo.Description, &todo.Date, &todo.Userid); err != nil {
		if err == sql.ErrNoRows {
			return todo, errors.New("No such row")
		}
		return todo, errors.New("Error getting todo")
	}
	return todo, nil
}
