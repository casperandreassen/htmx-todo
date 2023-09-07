package db

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

var DB *sqlx.DB

func InitDB() {
	f, err := os.OpenFile("todos.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	database, err := sqlx.Connect("sqlite3", "./todos.db")
	if err != nil {
		log.Fatalln(err)
	}

	// Create tables if not exists
	database.MustExec("create table if not exists user (id integer primary key not null, username text not null, password text not null); create table if not exists todos (id integer primary key not null, status integer not null, title text not null, description text not null, date datetime default null, userid integer not null, FOREIGN KEY(userid) REFERENCES user(id));")

	DB = database
}
