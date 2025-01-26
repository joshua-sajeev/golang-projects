package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

// DB is the global database instance
var DB *gorm.DB

// Connect initializes the database connection
func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm_app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful")
}
