package main

import (
	"Hunter-Hancock/dungeon-master/internal/server"
	"Hunter-Hancock/dungeon-master/pkg/ctrlc"
	"context"
	"log"
)

func main() {
	Server := server.NewServer()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrlc.HandleCtrlC(cancel)

	err := Server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server startup error: %v", err)
	}

}
