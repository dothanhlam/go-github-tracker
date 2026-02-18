package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Error writes a JSON error response
func Error(w http.ResponseWriter, status int, code, message string) {
	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Status:  status,
		},
	}
	JSON(w, status, resp)
}

// BadRequest writes a 400 Bad Request error
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized writes a 401 Unauthorized error
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// NotFound writes a 404 Not Found error
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

// InternalError writes a 500 Internal Server Error
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
