package httpapi

import "github.com/alisheikhattar/production-platform-lab/application/task-service/internal/task"

type createTaskRequest struct {
	Title string `json:"title"`
}

type updateTaskRequest struct {
	Status task.Status `json:"status"`
}
