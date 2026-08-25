package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskStore struct {
	db *pgxpool.Pool
}

func (s *TaskStore) Create(ctx context.Context, title string) (Task, error) {
	var task Task
	err := s.db.QueryRow(
		ctx,
		`INSERT INTO tasks (title)
		 VALUES ($1)
		 RETURNING id, title, completed`,
		title,
	).Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
func (s *TaskStore) UpdateCompleted(ctx context.Context, id int, completed bool) (Task, error) {
	var task Task
	err := s.db.QueryRow(ctx,
		`UPDATE tasks
						 SET completed = $1
						 WHERE id = $2
						 RETURNING title, id , completed`,
		completed,
		id,
	).Scan(&task.Title, &task.ID, &task.Completed)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (s *TaskStore) Delete(ctx context.Context, id int) error {
	result, err := s.db.Exec(
		ctx,
		`DELETE FROM tasks
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
func (s *TaskStore) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, title, completed
	 FROM tasks
	 ORDER BY id`)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return []Task{}, err
	}
	return tasks, nil
}
