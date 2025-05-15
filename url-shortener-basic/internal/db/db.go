package db

import (
	"database/sql"
	"log"
     _ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	connStr := "user=postgres password=123456 dbname=postgres sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Database!")
}
