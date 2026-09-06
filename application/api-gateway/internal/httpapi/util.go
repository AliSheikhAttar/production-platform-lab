package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func parseTaskID(r *http.Request) (int64, error) {
	return strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
}
