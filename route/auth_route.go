package route

import (
	"uas_be/app/model"
	"uas_be/app/service"
	"uas_be/helper"
	"uas_be/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterAuthRoutes mendaftarkan semua endpoint autentikasi
// Setiap route hanya: parse input → panggil service → return response
func RegisterAuthRoutes(app *fiber.App, authService service.AuthService) {
	auth := app.Group("/api/v1/auth")

	// POST /api/v1/auth/login
	auth.Post("/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
		}

		// Panggil service (semua business logic ada di sini)
		result, err := authService.Login(req.Username, req.Password)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusOK, "login berhasil", result)
	})

	// GET /api/v1/auth/profile
	auth.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		// Ambil user_id dari JWT context (dari middleware)
		userID := c.Locals("user_id").(string)

		// Panggil service untuk ambil profile
		profile, err := authService.GetUserProfile(userID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusOK, "profil user", profile)
	})
}
