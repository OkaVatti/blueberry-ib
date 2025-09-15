package main

import (
	"log"
	"os"

	"github.com/okavatti/blueberry/backend/server"
)

func main() {
	cfg := server.LoadConfigFromEnv()
	if err := server.Run(cfg); err != nil {
		log.Fatalf("server error: %v", err)
		os.Exit(1)
	}
}
