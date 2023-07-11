package main

import (
	"log"
	"paws-n-planes/internal/migrationManager/service/migrate"
)

func main() {
	migrationManager, err := migrate.New()

	if err != nil {
		log.Fatalf("Could not begin migration manager: %v", err)
	}

	migrationManager.Run()
}
