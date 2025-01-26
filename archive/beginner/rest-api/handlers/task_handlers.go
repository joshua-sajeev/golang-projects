package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api/models"
	"strconv"
	"strings"
)

var tasks = make(map[int]models.Task)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Decode the JSON body into a Task struct
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new ID based on the current number of tasks
	task.ID = len(tasks) + 1
	tasks[task.ID] = task

	// Log the tasks map for debugging
	fmt.Println("Tasks after creation:", tasks)

	// Respond with a status 201 (Created) and the task data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Return all tasks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from the URL
	// Assuming path format is /tasks/{id}
	// e.g., /tasks/1, /tasks/2
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse the task ID
	taskID := parts[2]

	// Convert the task ID to an integer
	id, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID format", http.StatusBadRequest)
		return
	}

	// Check if the task exists
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Return the task as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from the URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID := parts[2]
	id, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID format", http.StatusBadRequest)
		return
	}

	// Check if the task exists
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Decode the updated task from the request body
	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the task details
	task.Title = updatedTask.Title
	task.Completed = updatedTask.Completed
	tasks[id] = task

	// Return the updated task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from the URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskID := parts[2]
	id, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID format", http.StatusBadRequest)
		return
	}

	// Check if the task exists
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Delete the task from the map
	delete(tasks, id)

	// Return a success message
	w.WriteHeader(http.StatusNoContent)
}
