package models

import "net/http"

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// APIError represents an internal API error with context
type APIError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(message string, statusCode int, err error) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common error constructors
func NewBadRequestError(message string, err error) *APIError {
	return NewAPIError(message, http.StatusBadRequest, err)
}

func NewUnauthorizedError(message string, err error) *APIError {
	return NewAPIError(message, http.StatusUnauthorized, err)
}

func NewNotFoundError(message string, err error) *APIError {
	return NewAPIError(message, http.StatusNotFound, err)
}

func NewInternalServerError(message string, err error) *APIError {
	return NewAPIError(message, http.StatusInternalServerError, err)
}

// ToErrorResponse converts APIError to ErrorResponse
func (e *APIError) ToErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error:   http.StatusText(e.StatusCode),
		Message: e.Message,
		Code:    e.StatusCode,
	}
}
