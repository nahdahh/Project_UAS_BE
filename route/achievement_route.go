package route

import (
	"uas_be/app/service"
	"uas_be/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterAchievementRoutes mendaftarkan semua endpoint achievement
// Setiap route hanya: parse input → panggil service → return response
func RegisterAchievementRoutes(app *fiber.App, achievementService service.AchievementService) {
	achievements := app.Group("/api/v1/achievements")

	achievements.Use(middleware.AuthRequired())

	// GET /api/v1/achievements - List semua achievement
	achievements.Get("/", achievementService.HandleGetAll)

	// GET /api/v1/achievements/:id - Detail achievement milik student
	achievements.Get("/:id", achievementService.HandleGetByID)

	// POST /api/v1/achievements - Create achievement baru dengan status draft
	achievements.Post("/", achievementService.HandleCreate)

	// POST /api/v1/achievements/:id/submit - Submit prestasi untuk diverifikasi
	achievements.Post("/:id/submit", achievementService.HandleSubmit)

	// POST /api/v1/achievements/:id/verify - Verifikasi prestasi (dosen wali only)
	achievements.Post("/:id/verify", achievementService.HandleVerify)

	// POST /api/v1/achievements/:id/reject - Reject prestasi dengan alasan
	achievements.Post("/:id/reject", achievementService.HandleReject)

	// DELETE /api/v1/achievements/:id - Hapus prestasi (hanya draft)
	achievements.Delete("/:id", achievementService.HandleDelete)

	// GET /api/v1/achievements/advisee/list - List prestasi dari semua mahasiswa bimbingan
	achievements.Get("/advisee/list", achievementService.HandleGetAdviseeAchievements)
}
