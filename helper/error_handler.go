package helper

import "fmt"

// ApiError merepresentasikan standardisasi error dalam API
type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewApiError membuat ApiError baru
func NewApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

// Error implement error interface
func (e *ApiError) Error() string {
	return fmt.Sprintf("API Error [%d]: %s", e.Code, e.Message)
}
