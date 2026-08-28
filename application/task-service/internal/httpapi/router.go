package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/task"
)

type Handler struct {
	logger *slog.Logger
	tasks  task.Repository
}

func NewHandler(logger *slog.Logger, tasks task.Repository) *Handler {
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
	mux.HandleFunc("PATCH /tasks/{id}", h.updateTask)

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

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
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

	t, err := h.tasks.Create(
		r.Context(),
		request.Title,
	)
	if err != nil {
		h.logger.Error(
			"create task failed",
			slog.String("error", err.Error()),
		)

		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) listTasks(
	w http.ResponseWriter,
	r *http.Request,
) {
	tasks, err := h.tasks.List(r.Context())
	if err != nil {
		h.logger.Error(
			"list tasks failed",
			slog.String("error", err.Error()),
		)

		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
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

	t, err := h.tasks.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) updateTask(
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

	var request updateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if !validStatus(request.Status) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid task status",
		})
		return
	}

	t, err := h.tasks.UpdateStatus(
		r.Context(),
		id,
		request.Status,
	)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, t)
}
