package httpapi

type createTaskRequest struct {
	Title string `json:"title"`
}

type taskResponse struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
