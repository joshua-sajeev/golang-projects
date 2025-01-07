package models

import "gorm.io/gorm"

// Expense represents a spending record in the database.
type Category struct {
	gorm.Model
	Type string `json:"type"`
}
