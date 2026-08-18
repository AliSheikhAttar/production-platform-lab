package httpapi

import (
	"net/http"
	"strconv"

	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/task"
)

func parseTaskID(r *http.Request) (int64, error) {
	return strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
}

func validStatus(status task.Status) bool {
	switch status {
	case task.StatusPending,
		task.StatusProcessing,
		task.StatusCompleted,
		task.StatusFailed:
		return true

	default:
		return false
	}
}
