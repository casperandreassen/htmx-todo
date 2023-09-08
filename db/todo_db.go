package db

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-server/utils"
	"log"
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

func InsertTodo(todo Todo, userid int) (int64, error) {
	result, err := DB.Exec("INSERT INTO todos (title, description, status, date, userid) VALUES (?, ?, ?, ?, ?)", todo.Title, todo.Description, 0, utils.NewNullString(todo.Date), userid)
	if err != nil {
		return 0, errors.New("Could not insert todo")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("Could not get id of inserted row")
	}
	return id, nil
}

func DeleteTodo(todoId int) error {
	_, err := DB.Exec("DELETE FROM todos WHERE id = ?", todoId)
	if err != nil {
		return errors.New("Could not delete todo")
	}
	return nil
}

func UpdateTodoStatus(todoId int, newStatus int) error {
	_, err := DB.Exec("UPDATE todos SET status = ? where id = ?", newStatus, todoId)
	if err != nil {
		return errors.New("Could not update todo")
	}
	return nil
}

func GetSingleTodo(todoId int) (DBTodo, error) {
	var todo DBTodo

	row := DB.QueryRow("SELECT * FROM todos WHERE todos.id = ?", todoId)
	if err := row.Scan(&todo.Id, &todo.Status, &todo.Title, &todo.Description, &todo.Date, &todo.Userid); err != nil {
		if err == sql.ErrNoRows {
			return todo, errors.New("No such row")
		}
		return todo, errors.New("Error getting todo")
	}
	return todo, nil
}

func GetTodos(c *gin.Context) ([]DBTodo, error) {
	userid := c.MustGet("id").(float64)
	todos := []DBTodo{}
	query, err := DB.Prepare("SELECT * FROM todos WHERE userid = ?")

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

func TransformTodos(todos []DBTodo) ([]DBTodo, []DBTodo, []DBTodo, error) {
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
