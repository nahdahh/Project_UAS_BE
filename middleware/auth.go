package middleware

import (
	"strings"
	"uas_be/helper"
	"uas_be/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthRequired middleware untuk validasi JWT token dan inject user context
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "token tidak ditemukan")
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "format token tidak valid")
		}

		tokenString := parts[1]

		// Verify dan parse JWT token
		claims, err := utils.GetClaimsFromToken(tokenString)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "token tidak valid: "+err.Error())
		}

		// Inject user info ke context agar bisa diakses di route handler
		c.Locals("user_id", claims.Sub)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

// AdminOnly middleware untuk check apakah user adalah admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Pastikan AuthRequired() middleware sudah di-apply terlebih dahulu
		role, ok := c.Locals("role").(string)
		if !ok || role != "Admin" {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "hanya admin yang dapat mengakses resource ini")
		}

		return c.Next()
	}
}

// DosenWaliOnly middleware untuk check apakah user adalah dosen wali
func DosenWaliOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "Dosen Wali" {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "hanya dosen wali yang dapat mengakses resource ini")
		}

		return c.Next()
	}
}

// MahasiswaOnly middleware untuk check apakah user adalah mahasiswa
func MahasiswaOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "Mahasiswa" {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "hanya mahasiswa yang dapat mengakses resource ini")
		}

		return c.Next()
	}
}

// HasPermission middleware untuk check apakah user memiliki specific permission
func HasPermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil permissions dari context (set oleh AuthRequired middleware)
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "permissions tidak ditemukan")
		}

		// Check apakah required permission ada di dalam permissions list
		hasPermission := false
		for _, perm := range permissions {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "anda tidak memiliki permission: "+requiredPermission)
		}

		return c.Next()
	}
}

// OptionalAuth middleware untuk optional JWT validation (tidak required tapi bisa ada)
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Token tidak ada, tapi kita izinkan
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		// Try verify token, tapi jangan fail jika invalid
		claims, err := utils.GetClaimsFromToken(tokenString)
		if err == nil {
			// Inject ke context jika valid
			c.Locals("user_id", claims.Sub)
			c.Locals("role", claims.Role)
		}

		return c.Next()
	}
}
