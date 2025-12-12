package middleware

import (
	"strings"
	"uas_be/helper"
	"uas_be/utils"

	"github.com/gofiber/fiber/v2"
)

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
		c.Locals("userID", claims.Sub)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

func RBACMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "permissions tidak ditemukan")
		}

		// Check if user has the required permission
		for _, perm := range permissions {
			if perm == requiredPermission {
				return c.Next()
			}
		}

		return helper.ErrorResponse(c, fiber.StatusForbidden, "anda tidak memiliki permission: "+requiredPermission)
	}
}

func AuthMiddleware() fiber.Handler {
	return AuthRequired()
}

func HasPermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "permissions tidak ditemukan")
		}

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

func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		claims, err := utils.GetClaimsFromToken(tokenString)
		if err == nil {
			c.Locals("userID", claims.Sub)
			c.Locals("role", claims.Role)
		}

		return c.Next()
	}
}
