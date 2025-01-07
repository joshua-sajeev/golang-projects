package models

import "gorm.io/gorm"

// Expense represents a spending record in the database.
type Expense struct {
	gorm.Model
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}
