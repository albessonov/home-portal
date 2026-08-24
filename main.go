package main

import (
	"fmt"
	"net/http"
)

func main() {
	app := App{
		tasks: nil,
	}
	http.HandleFunc("GET /api/dashboard", (app.dashboardHandler))
	http.HandleFunc("GET /api/health", (healthHandler))
	http.HandleFunc("POST /api/tasks", app.addTaskHandler)
	http.HandleFunc("PATCH /api/tasks/{id}", app.updateTaskHandler)
	http.HandleFunc("DELETE /api/tasks/{id}", app.deleteTaskHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
