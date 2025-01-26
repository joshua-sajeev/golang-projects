package main

import (
	"github.com/gin-gonic/gin"
	"todo-api/routes"
	"todo-api/storage"
)

func main() {
	storage.InitDb()
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8081")
}
