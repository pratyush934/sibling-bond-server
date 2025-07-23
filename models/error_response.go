package models

import "fmt"

type ErrorResponse struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	InternalError error  `json:"internalError"`
}

type HTTPError struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	InternalError error  `json:"-"`
}

func (h *HTTPError) Error() string {
	if h.InternalError != nil {
		return fmt.Sprintf("There is an issue %v, %v", h.InternalError, h.Message)
	}
	return h.Message
}

func NewError(status int, message string, internalError error) *HTTPError {
	return &HTTPError{
		Status:        status,
		Message:       message,
		InternalError: internalError,
	}
}
