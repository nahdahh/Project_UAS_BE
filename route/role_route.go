package route

import (
	"uas_be/app/service"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoleRoutes mendaftarkan semua endpoint untuk role management (admin only)
// Route hanya: parse input → panggil service → return response
func RegisterRoleRoutes(app *fiber.App, roleService service.RoleService) {
	roles := app.Group("/api/v1/roles")

	// GET /api/v1/roles - List semua role
	roles.Get("/", roleService.HandleGetAllRoles)

	// GET /api/v1/roles/:id - Detail role
	roles.Get("/:id", roleService.HandleGetRoleByID)

	// POST /api/v1/roles - Buat role baru
	roles.Post("/", roleService.HandleCreateRole)

	// PUT /api/v1/roles/:id - Update role
	roles.Put("/:id", roleService.HandleUpdateRole)

	// POST /api/v1/roles/:id/permissions/:permissionId - Assign permission ke role
	roles.Post("/:id/permissions/:permissionId", roleService.HandleAssignPermissionToRole)

	// DELETE /api/v1/roles/:id/permissions/:permissionId - Remove permission dari role
	roles.Delete("/:id/permissions/:permissionId", roleService.HandleRemovePermissionFromRole)
}
