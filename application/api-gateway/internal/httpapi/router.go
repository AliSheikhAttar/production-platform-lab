package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	taskv1 "github.com/alisheikhattar/production-platform-lab/application/task-service/gen/task/v1"
)

type Handler struct {
	logger *slog.Logger
	tasks  taskv1.TaskServiceClient
}

func NewHandler(
	logger *slog.Logger,
	tasks taskv1.TaskServiceClient,
) *Handler {
	return &Handler{
		logger: logger,
		tasks:  tasks,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health/live", h.liveness)
	mux.HandleFunc("GET /health/ready", h.readiness)
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("GET /tasks/{id}", h.getTask)

	return mux
}

func (h *Handler) liveness(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "alive",
	})
}

func (h *Handler) readiness(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func (h *Handler) createTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	request.Title = strings.TrimSpace(request.Title)
	if request.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
		return
	}

	resp, err := h.tasks.CreateTask(
		r.Context(),
		&taskv1.CreateTaskRequest{
			Title: request.Title,
		},
	)
	if err != nil {
		h.handleGRPCError(w, err, "create task failed")
		return
	}

	writeJSON(w, http.StatusCreated, toTaskResponse(resp.GetTask()))
}

func (h *Handler) listTasks(
	w http.ResponseWriter,
	r *http.Request,
) {
	resp, err := h.tasks.ListTasks(
		r.Context(),
		&taskv1.ListTasksRequest{},
	)
	if err != nil {
		h.handleGRPCError(w, err, "list tasks failed")
		return
	}

	tasks := make([]taskResponse, 0, len(resp.GetTasks()))
	for _, t := range resp.GetTasks() {
		tasks = append(tasks, toTaskResponse(t))
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) getTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := parseTaskID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid task id",
		})
		return
	}

	resp, err := h.tasks.GetTask(
		r.Context(),
		&taskv1.GetTaskRequest{
			Id: id,
		},
	)
	if err != nil {
		h.handleGRPCError(w, err, "get task failed")
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(resp.GetTask()))
}

func (h *Handler) handleGRPCError(
	w http.ResponseWriter,
	err error,
	logMessage string,
) {
	h.logger.Error(
		logMessage,
		slog.String("error", err.Error()),
	)

	status, message := grpcStatusToHTTP(err)
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func toTaskResponse(t *taskv1.Task) taskResponse {
	return taskResponse{
		ID:     t.GetId(),
		Title:  t.GetTitle(),
		Status: t.GetStatus(),
	}
}
