package helper

import (
	"strconv"

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
