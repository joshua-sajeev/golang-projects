package routes

import (
	"todo-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": true})
	})
	r.GET("/todos", handlers.GetTodos)
	r.POST("/todos", handlers.CreateTodo)
}
