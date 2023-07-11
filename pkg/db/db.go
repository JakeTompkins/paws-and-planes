package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var lock = &sync.Mutex{}

func GetInstance() (*sql.DB, error) {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()

		db, err := sql.Open("mysql", os.Getenv("DSN"))

		if err != nil {
			log.Printf("failed to connect: %v", err)

			return nil, err
		}

		if err := db.Ping(); err != nil {
			log.Printf("failed to ping: %v", err)

			return nil, err
		}

		log.Println("Successfully connected to PlanetScale!")
	}

	return db, nil
}
