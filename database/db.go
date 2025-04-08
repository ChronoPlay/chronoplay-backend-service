package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// dsn := os.Getenv("MYSQL_DSN") // Example: "user:password@tcp(localhost:3306)/yourdb?parseTime=true"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Failed to connect to DB:", err)
	// }

	// DB = db
}
