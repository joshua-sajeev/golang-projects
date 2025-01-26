package storage

import (
	"log"
	"todo-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB // Add this global variable

func InitDb() {
	var err error
	DB, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{}) // Changed to todos.db
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Running database migrations...")
	DB.AutoMigrate(&models.Todo{})
	log.Println("Database initialized.")
}
