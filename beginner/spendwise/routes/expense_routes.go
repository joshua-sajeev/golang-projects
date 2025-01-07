package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterExpenseRoutes sets up routes for expense operations.
func RegisterExpenseRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/expenses", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "List all expenses"})
		})
		api.POST("/expenses", func(c *gin.Context) {
			c.JSON(201, gin.H{"message": "Create an expense"})
		})
	}
}
