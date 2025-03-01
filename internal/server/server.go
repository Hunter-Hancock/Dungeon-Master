package server

import (
	"Hunter-Hancock/dungeon-master/internal/db"
	"database/sql"
	"net/http"
	"time"
)

type Server struct {
	db *sql.DB
}

func NewServer() *http.Server {
	db, err := db.OpenDB()
	if err != nil {
		panic(err)
	}
	newServer := &Server{
		db: db,
	}

	server := &http.Server{
		Addr:         ":6969",
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
