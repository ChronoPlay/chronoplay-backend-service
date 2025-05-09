package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	var db *gorm.DB
	var err error

	maxRetries := 10
	for retries := 0; retries < maxRetries; retries++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to the database successfully.")
			DB = db
			return
		}

		log.Printf("Retry %d/%d: Failed to connect to DB: %v", retries+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to the database after several attempts:", err)
}
