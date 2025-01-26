package models

type Todo struct {
	ID        string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Title     string `gorm:"not null" json:"title" validate:"required"`
	Completed bool   `json:"completed"`
}
