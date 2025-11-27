package route

import (
	"uas_be/app/model"
	"uas_be/app/service"
	"uas_be/helper"
	"uas_be/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes register user management endpoints
func RegisterUserRoutes(app *fiber.App, userService service.UserService) {
	users := app.Group("/api/v1/users")

	// GET /api/v1/users - admin only
	users.Get("/", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		page, pageSize := helper.ExtractPaginationParams(c)

		//ambil semua user dari service
		users, total, err := userService.GetAllUsers(page, pageSize)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusOK, "data user berhasil diambil", fiber.Map{
			"data":  users,
			"total": total,
			"page":  page,
		})
	})

	// GET /api/v1/users/:id
	users.Get("/:id", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		userID := c.Params("id")

		//ambil user dari service
		user, err := userService.GetUserByID(userID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusOK, "user berhasil diambil", user)
	})

	// POST /api/v1/users
	users.Post("/", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		var req model.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
		}

		//buat user baru di service
		user, err := userService.CreateUser(&req, req.RoleID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusCreated, "user berhasil dibuat", user)
	})

	// PUT /api/v1/users/:id
	users.Put("/:id", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			FullName string `json:"full_name"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
		}

		//update user di service
		err := userService.UpdateUser(c.Params("id"), req.Username, req.Email, req.FullName)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusOK, "user berhasil diupdate", nil)
	})

	// DELETE /api/v1/users/:id
	users.Delete("/:id", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		userID := c.Params("id")

		//hapus user di service
		err := userService.DeleteUser(userID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		return helper.SuccessResponse(c, fiber.StatusOK, "user berhasil dihapus", nil)
	})

	// PUT /api/v1/users/:id/role
	users.Put("/:id/role", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		var req struct {
			RoleID string `json:"role_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
		}

		//TODO: implement assign role di service
		return helper.ErrorResponse(c, fiber.StatusNotImplemented, "fitur belum diimplementasikan")
	})
}
