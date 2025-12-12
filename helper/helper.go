package helper

import (
	"strconv"
	"uas_be/app/model"

	"github.com/gofiber/fiber/v2"
)

// ResponseFormat format response API standar
type ResponseFormat struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse kirim response sukses
func SuccessResponse(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(ResponseFormat{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse kirim response error
func ErrorResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(ResponseFormat{
		Status:  "error",
		Message: message,
		Error:   message,
	})
}

// ExtractPaginationParams extract pagination dari query parameter
func ExtractPaginationParams(c *fiber.Ctx) (int, int) {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	return page, pageSize
}

func GetStatusCodeFromError(err error) int {
	return getStatusCodeFromError(err)
}

// getStatusCodeFromError menentukan status code HTTP dari error message
func getStatusCodeFromError(err error) int {
	errMsg := err.Error()

	// Not found errors
	if contains(errMsg, "tidak ditemukan") || contains(errMsg, "not found") {
		return fiber.StatusNotFound
	}

	// Unauthorized/Forbidden errors
	if contains(errMsg, "unauthorized") || contains(errMsg, "bukan milik anda") ||
		contains(errMsg, "bukan advisor") {
		return fiber.StatusForbidden
	}

	// Conflict errors
	if contains(errMsg, "sudah terdaftar") || contains(errMsg, "sudah ada") ||
		contains(errMsg, "already exists") {
		return fiber.StatusConflict
	}

	// Bad request errors
	if contains(errMsg, "tidak boleh kosong") || contains(errMsg, "tidak valid") ||
		contains(errMsg, "hanya") || contains(errMsg, "harus") {
		return fiber.StatusBadRequest
	}

	// Default to internal server error
	return fiber.StatusInternalServerError
}

// contains helper untuk check substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(str) > len(substr) &&
		(hasSubstring(str, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SuccessResponse returns success response without needing fiber.Ctx
func SuccessResponseSimple(message string, data interface{}) model.APIResponse {
	return model.APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// ErrorResponse returns error response without needing fiber.Ctx
func ErrorResponseSimple(message string) model.APIResponse {
	return model.APIResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}
