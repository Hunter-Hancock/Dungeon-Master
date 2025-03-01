package db

import (
	"database/sql"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"log"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("libsql", "../test.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
