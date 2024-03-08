package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Task represents a to-do task.
type Task struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"` // e.g., "pending" or "completed"
}

var tasks = []Task{} // In-memory storage for tasks
const Dport = ":8012"
const (
	studentID       = "500230292"
	githubRepoURL   = "https://github.com/hammadCodes/goCRUDDocker"
	apiGuideMessage = `
		This is a simple to-do task server. Which is integared with Docker as well. This server allows you to perform CRUD (Create, Read, Update, Delete) operations on tasks.

		- To view this guide, send a GET request to /.
		- To view all tasks, send a GET request to /tasks.
		- To create a new task, send a POST request to /tasks with a JSON body containing "title" and "status" fields.
			Example: {"title": "Task 1", "status": "pending"}
		- To update a task, send a PUT request to /task/{id} with the updated task details in the JSON body.
		- To delete a task, send a DELETE request to /task/{id}.
		- A new root handler function has been added to the server to display the API guide message.

		The more documentation for github at: https://github.com/hammadCodes/goCRUDDocker
		The more documentation for docker at: https://hub.docker.com/repository/docker/hammadcodes/gocruddocker/general

		Student ID: %s
		GitHub Repository: %s
		`
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/task/", taskHandler)
	fmt.Printf("Server is starting on port: %v\n", Dport) // Added newline for better terminal output
	http.ListenAndServe(Dport, nil)
}

// Handle requests to the "/" endpoint
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, apiGuideMessage, studentID, githubRepoURL)
}

// Handle requests to the /tasks endpoint
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(tasks)
	case "POST":
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.ID = uuid.New().String() // Generate a unique ID for the task
		tasks = append(tasks, task)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handle requests to the /task/{id} endpoint
func taskHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the task ID from the URL path
	taskID := strings.TrimPrefix(r.URL.Path, "/task/")

	switch r.Method {
	case "PUT":
		var updatedTask Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		found := false
		for i, task := range tasks {
			if task.ID == taskID {
				updatedTask.ID = task.ID // Ensure the ID remains unchanged
				tasks[i] = updatedTask
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(updatedTask)
	case "DELETE":
		index := -1
		for i, task := range tasks {
			if task.ID == taskID {
				index = i
				break
			}
		}
		if index != -1 {
			tasks = append(tasks[:index], tasks[index+1:]...)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Task not found", http.StatusNotFound)
		}
	}
}
