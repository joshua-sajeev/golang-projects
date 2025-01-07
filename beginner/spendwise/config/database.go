package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupDatabase initializes the SQLite database and returns the GORM DB instance.
func SetupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("spendwise.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}
