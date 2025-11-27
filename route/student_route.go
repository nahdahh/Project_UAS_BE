package route

import (
	"github.com/gofiber/fiber/v2"
	"uas_be/app/service"
	"uas_be/middleware"
)

// RegisterStudentRoutes mendaftarkan semua endpoint student
// Route hanya setup endpoint, service Handle yang urus semuanya
func RegisterStudentRoutes(app *fiber.App, studentService service.StudentService) {
	students := app.Group("/api/v1/students")

	// GET semua students
	students.Get("/", middleware.AuthRequired(), studentService.HandleGetAll)

	// GET student by ID
	students.Get("/:id", middleware.AuthRequired(), studentService.HandleGetByID)

	// POST create student
	students.Post("/", middleware.AuthRequired(), studentService.HandleCreate)

	// PUT update student
	students.Put("/:id", middleware.AuthRequired(), studentService.HandleUpdate)

	// DELETE student
	students.Delete("/:id", middleware.AuthRequired(), studentService.HandleDelete)
}
