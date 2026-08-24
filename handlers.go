package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	a := HealthResponse{
		Status: "ok",
	}
	json.NewEncoder(w).Encode(a)

}
func (a *App) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := DashboardResponse{
		Greetings: "Добрый вечер",
		Tasks:     a.tasks,
		Shopping:  []string{},
	}
	json.NewEncoder(w).Encode(d)
}
func (a *App) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	newTask := Task{
		ID:        len(a.tasks) + 1,
		Title:     req.Title,
		Completed: false,
	}
	a.tasks = append(a.tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}
func (a *App) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	var req UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Completed == nil {
		http.Error(w, "hueta", http.StatusBadRequest)
		return
	}
	for i := range a.tasks {
		if a.tasks[i].ID == id {
			a.tasks[i].Completed = *req.Completed
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(a.tasks[i])
			return
		}
	}
	http.Error(w, "id not found", http.StatusNotFound)
	fmt.Println(id)
}
