package route

import (
	"uas_be/app/service"
	"uas_be/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App, authService service.AuthService) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.RefreshToken)
	auth.Post("/logout", middleware.AuthRequired(), authService.Logout)
	auth.Get("/profile", middleware.AuthRequired(), authService.GetProfile)
}
