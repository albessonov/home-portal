package main

type HealthResponse struct {
	Status string `json:"status"`
}

type DashboardResponse struct {
	Greetings string   `json:"greeting"`
	Tasks     []Task   `json:"tasks"`
	Shopping  []string `json:"shopping"`
}
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
type CreateTaskRequest struct {
	Title string `json:"title"`
}
type UpdateTaskRequest struct {
	Completed *bool `json:"completed"`
}
