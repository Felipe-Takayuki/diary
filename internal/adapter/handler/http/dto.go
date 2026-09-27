package http

// ItemDTO represents the JSON payload format for an individual goal.
type ItemDTO struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// MetasResponse represents the response payload for GET /api/metas.
type MetasResponse struct {
	Date  string    `json:"date"`
	Items []ItemDTO `json:"items"`
}

// SaveMetasRequest represents the request payload for POST /api/metas.
type SaveMetasRequest struct {
	Date  string    `json:"date"`
	Items []ItemDTO `json:"items"`
}

// StatusResponse represents simple operational responses.
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents standard JSON error format.
type ErrorResponse struct {
	Error string `json:"error"`
}
