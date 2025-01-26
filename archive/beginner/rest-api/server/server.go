package server

import (
	"log"
	"net/http"
	"rest-api/handlers"
)

type Server struct {
	Addr   string
	Router *http.ServeMux
}

// Create a new server
func NewServer(addr string) *Server {
	server := http.NewServeMux()
	server.HandleFunc("GET /", handlers.GetAllTasks)
	server.HandleFunc("GET /tasks", handlers.GetAllTasks)
	server.HandleFunc("GET /tasks/{id}", handlers.GetTask)
	server.HandleFunc("PUT /tasks/{id}", handlers.UpdateTask)
	server.HandleFunc("DELETE /tasks/{id}", handlers.DeleteTask)
	server.HandleFunc("POST /tasks", handlers.CreateTask)

	return &Server{
		Addr:   addr,
		Router: server,
	}
}
func (s *Server) Run() error {

	log.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(s.Addr, s.Router); err != nil {
		log.Fatalf("Server failed: %v", err)
		return err
	}
	return nil
}
