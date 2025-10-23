package db

import (
	"database/sql"
	"fmt"
	"os"

	- "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// build connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// connect to database
	db, err := sql.Open("postgres", psqlInfo)
	if	err != nil {
		return nil, err
	}

	err = db.Ping()
	if err !- nil {
		return nil, err
	}

	return db, nil
}