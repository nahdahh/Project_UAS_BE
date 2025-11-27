package route

import (
	"uas_be/app/service"
	"uas_be/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterLecturerRoutes mendaftarkan semua route untuk lecturer/dosen
func RegisterLecturerRoutes(app *fiber.App, lecturerService service.LecturerService) {
	lecturers := app.Group("/api/v1/lecturers")

	lecturers.Use(middleware.AuthRequired())

	// GET /api/v1/lecturers - Ambil semua lecturer dengan pagination
	lecturers.Get("/", lecturerService.HandleGetAll)

	// GET /api/v1/lecturers/:id - Ambil lecturer berdasarkan ID
	lecturers.Get("/:id", lecturerService.HandleGetByID)

	// GET /api/v1/lecturers/:id/advisees - Ambil daftar mahasiswa bimbingan
	lecturers.Get("/:id/advisees", lecturerService.HandleGetAdvisees)

	// POST /api/v1/lecturers - Buat lecturer baru
	lecturers.Post("/", lecturerService.HandleCreate)

	// PUT /api/v1/lecturers/:id - Update lecturer
	lecturers.Put("/:id", lecturerService.HandleUpdate)

	// DELETE /api/v1/lecturers/:id - Hapus lecturer
	lecturers.Delete("/:id", lecturerService.HandleDelete)
}
