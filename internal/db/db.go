package db

import (
	"database/sql"
	"fmt"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"log"
	"os"
)

func OpenDB() (*sql.DB, error) {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to database")

	return db, nil
}
