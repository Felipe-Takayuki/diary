package http

// ItemDTO represents the JSON payload format for an individual goal.
type ItemDTO struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// MetasResponse represents the response payload for GET /api/metas and GET /api/goals.
type MetasResponse struct {
	Date  string    `json:"date"`
	Items []ItemDTO `json:"items"`
}

// GoalsResponse is an alias for MetasResponse.
type GoalsResponse = MetasResponse

// SaveMetasRequest represents the request payload for POST /api/metas and POST /api/goals.
type SaveMetasRequest struct {
	Date  string    `json:"date"`
	Items []ItemDTO `json:"items"`
}

// SaveGoalsRequest is an alias for SaveMetasRequest.
type SaveGoalsRequest = SaveMetasRequest

// StatusResponse represents simple operational responses.
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents standard JSON error format.
type ErrorResponse struct {
	Error string `json:"error"`
}
