package main

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title string
	Slug  string
}

var db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

func main() {

	db.AutoMigrate(&Post{})

}
