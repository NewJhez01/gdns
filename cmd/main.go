package main

import (
	"log"

	"github.com/gdns/cmd/server"
)

func main() {
	err := server.Serve()
	if err != nil {
		log.Fatalf("failed to start server prev: %s", err)
	}
}
