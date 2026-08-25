package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://home:home@localhost:5432/home_portal",
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		fmt.Println(err)
		return
	}

	taskStore := &TaskStore{
		db: db,
	}

	app := App{
		tasks: taskStore,
	}

	http.HandleFunc("GET /api/dashboard", app.dashboardHandler)
	http.HandleFunc("GET /api/health", healthHandler)
	http.HandleFunc("POST /api/tasks", app.addTaskHandler)
	http.HandleFunc("PATCH /api/tasks/{id}", app.updateTaskHandler)
	http.HandleFunc("DELETE /api/tasks/{id}", app.deleteTaskHandler)

	fs := http.FileServer(http.Dir("./web"))

	http.Handle(
		"GET /static/",
		http.StripPrefix("/static/", fs),
	)

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
