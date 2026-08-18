package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/task"
)

func TestLiveness(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	repository := task.NewMemoryRepository()

	handler := NewHandler(
		logger,
		repository,
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/health/live",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}
