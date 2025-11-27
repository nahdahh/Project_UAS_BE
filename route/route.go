package route

import (
	"uas_be/app/service"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes mendaftarkan semua API routes sesuai SRS
// Struktur: /api/v1/{resource}/{action}
func RegisterRoutes(
	app *fiber.App,
	authService service.AuthService,
	userService service.UserService,
	studentService service.StudentService,
	lecturerService service.LecturerService,
	achievementService service.AchievementService,
	roleService service.RoleService,
) {
	// Route untuk autentikasi - FR-001, FR-002
	RegisterAuthRoutes(app, authService)

	// Route untuk manajemen role dan permission (admin only)
	RegisterRoleRoutes(app, roleService)

	// Route untuk manajemen user - FR-009 (admin only)
	RegisterUserRoutes(app, userService)

	// Route untuk manajemen student
	RegisterStudentRoutes(app, studentService)

	// Route untuk manajemen lecturer - FR-006
	RegisterLecturerRoutes(app, lecturerService)

	// Route untuk manajemen achievement - FR-003 sampai FR-011
	RegisterAchievementRoutes(app, achievementService)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0",
			"service": "UAS_BE - Sistem Pelaporan Prestasi Mahasiswa",
		})
	})
}
