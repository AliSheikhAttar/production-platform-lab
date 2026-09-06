package grpcapi

import (
	"context"
	"errors"
	"strings"

	taskv1 "github.com/alisheikhattar/production-platform-lab/application/task-service/gen/task/v1"
	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	taskv1.UnimplementedTaskServiceServer

	tasks task.Repository
}

func NewServer(tasks task.Repository) *Server {
	return &Server{
		tasks: tasks,
	}
}

func (s *Server) CreateTask(
	ctx context.Context,
	req *taskv1.CreateTaskRequest,
) (*taskv1.CreateTaskResponse, error) {
	title := strings.TrimSpace(req.GetTitle())
	if title == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"title is required",
		)
	}

	t, err := s.tasks.Create(ctx, title)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"create task: %v",
			err,
		)
	}

	return &taskv1.CreateTaskResponse{
		Task: toProtoTask(t),
	}, nil
}

func (s *Server) GetTask(
	ctx context.Context,
	req *taskv1.GetTaskRequest,
) (*taskv1.GetTaskResponse, error) {
	t, err := s.tasks.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"task not found",
			)
		}

		return nil, status.Errorf(
			codes.Internal,
			"get task: %v",
			err,
		)
	}

	return &taskv1.GetTaskResponse{
		Task: toProtoTask(t),
	}, nil
}

func (s *Server) ListTasks(
	ctx context.Context,
	_ *taskv1.ListTasksRequest,
) (*taskv1.ListTasksResponse, error) {
	tasks, err := s.tasks.List(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"list tasks: %v",
			err,
		)
	}

	protoTasks := make([]*taskv1.Task, 0, len(tasks))
	for _, t := range tasks {
		protoTasks = append(protoTasks, toProtoTask(t))
	}

	return &taskv1.ListTasksResponse{
		Tasks: protoTasks,
	}, nil
}

func toProtoTask(t task.Task) *taskv1.Task {
	return &taskv1.Task{
		Id:     t.ID,
		Title:  t.Title,
		Status: string(t.Status),
	}
}
