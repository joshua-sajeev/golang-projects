package main

import (
	"fmt"
	"rest-api/server"
)

func main() {
	server := server.NewServer(":8080")
	server.Run()
	fmt.Println("Hello world!")
}
