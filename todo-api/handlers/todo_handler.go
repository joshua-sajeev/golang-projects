package handlers

import (
	"github.com/gin-gonic/gin"
	"todo-api/models"
	"todo-api/storage"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo

	// Fetch all todos from database
	result := storage.DB.Find(&todos)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Could not fetch todos"})
		return
	}

	if len(todos) == 0 {
		c.JSON(200, gin.H{
			"message": "🎉 No todos found! Time to binge-watch Netflix?",
			"data":    todos,
		})
		return
	}

	c.JSON(200, gin.H{"data": todos})
}

func CreateTodo(c *gin.Context) {
	var newTodo models.Todo

	// Parse JSON input
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Create todo in database
	result := storage.DB.Create(&newTodo)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Could not create todo"})
		return
	}

	c.JSON(201, newTodo)
}
