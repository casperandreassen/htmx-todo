package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
	"log"
	"os"
)

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

var DB *sql.DB

func InitDB() {
	env := os.Getenv("APP_ENV")

	if env == "DEV" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatalf("Some error occured. Err: %s", err)
		}
	}

	turso_auth_key := os.Getenv("TURSO_AUTH_KEY")

	if turso_auth_key == "" {
		log.Fatalf("TURSO AUTH KEY not present.")
	}

	database, err := sql.Open("libsql", "libsql://pleased-the-leader-casperandreassen.turso.io?authToken="+turso_auth_key)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	pingErr := database.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// Create tables if not exists
	database.Exec("create table if not exists user (id integer primary key not null, username text unique not null, password text not null); create table if not exists todos (id integer primary key not null, status integer not null, title text not null, description text not null, date datetime default null, userid integer not null, FOREIGN KEY(userid) REFERENCES user(id));")

	fmt.Println("Database successfully connected!")

	DB = database

}
