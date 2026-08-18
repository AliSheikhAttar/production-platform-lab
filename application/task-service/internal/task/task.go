package task

import (
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
	Create(title string) Task
	Get(id int64) (Task, error)
	List() []Task
	UpdateStatus(id int64, status Status) (Task, error)
}
