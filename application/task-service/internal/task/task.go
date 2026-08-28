package task

import (
	"context"
	"errors"
	"time"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(
		ctx context.Context,
		title string,
	) (Task, error)

	Get(
		ctx context.Context,
		id int64,
	) (Task, error)

	List(
		ctx context.Context,
	) ([]Task, error)

	UpdateStatus(
		ctx context.Context,
		id int64,
		status Status,
	) (Task, error)
}
