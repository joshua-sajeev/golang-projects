package main

import (
	"github.com/gin-gonic/gin"
	"spendwise/config"
	"spendwise/routes"
)

func main() {
	// Initialize database
	db := config.SetupDatabase()

	// Close database connection when the app shuts down
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	routes.RegisterExpenseRoutes(router)

	// Start server
	router.Run(":8080")
}
