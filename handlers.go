package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	a := HealthResponse{
		Status: "ok",
	}
	json.NewEncoder(w).Encode(a)

}
func (a *App) dashboardHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := a.tasks.List(r.Context())
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	d := DashboardResponse{
		Greetings: "Добрый вечер",
		Tasks:     tasks,
		Shopping:  []string{},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
func (a *App) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	task, err := a.tasks.Create(
		r.Context(), req.Title,
	)
	if err != nil {
		http.Error(w, "task was not created", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

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

	task, err := a.tasks.UpdateCompleted(r.Context(), id, *req.Completed)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	err = a.tasks.Delete(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
