package httpapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	taskv1 "github.com/alisheikhattar/production-platform-lab/application/task-service/gen/task/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeTaskClient struct{}

func (fakeTaskClient) CreateTask(
	_ context.Context,
	req *taskv1.CreateTaskRequest,
	_ ...grpc.CallOption,
) (*taskv1.CreateTaskResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	return &taskv1.CreateTaskResponse{
		Task: &taskv1.Task{
			Id:     1,
			Title:  req.GetTitle(),
			Status: "pending",
		},
	}, nil
}

func (fakeTaskClient) GetTask(
	context.Context,
	*taskv1.GetTaskRequest,
	...grpc.CallOption,
) (*taskv1.GetTaskResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (fakeTaskClient) ListTasks(
	context.Context,
	*taskv1.ListTasksRequest,
	...grpc.CallOption,
) (*taskv1.ListTasksResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func TestCreateTask(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	handler := NewHandler(logger, fakeTaskClient{})

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title":"learn docker"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}
}

func TestLiveness(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	handler := NewHandler(logger, fakeTaskClient{})

	request := httptest.NewRequest(
		http.MethodGet,
		"/health/live",
		nil,
	)

	recorder := httptest.NewRecorder()
	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}
