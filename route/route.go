package route

import (
	"uas_be/app/repository"
	"uas_be/app/service"
	"uas_be/database"
	"uas_be/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	db := database.GetDB()

	achievementRepo := repository.NewAchievementRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo, permissionRepo)
	achievementService := service.NewAchievementService(achievementRepo, studentRepo)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo)
	studentService := service.NewStudentService(studentRepo, achievementRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	userService := service.NewUserService(userRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	reportService := service.NewReportService(achievementRepo, studentRepo)

	SetupAuthRoutes(app, authService)
	SetupAchievementRoutes(app, achievementService)
	SetupLecturerRoutes(app, lecturerService)
	SetupStudentRoutes(app, studentService)
	SetupRoleRoutes(app, roleService)
	SetupUserRoutes(app, userService)
	SetupPermissionRoutes(app, permissionService)
	SetupReportRoutes(app, reportService)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0",
			"service": "UAS_BE - Sistem Pelaporan Prestasi Mahasiswa",
		})
	})
}

func SetupAchievementRoutes(app *fiber.App, achievementService service.AchievementService) {
	group := app.Group("/api/v1/achievements", middleware.AuthMiddleware())

	group.Get("/", middleware.RBACMiddleware("achievement:read"), achievementService.GetAllAchievements)
	group.Get("/:id", middleware.RBACMiddleware("achievement:read"), achievementService.GetAchievementDetail)
	group.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), achievementService.GetAchievementHistory)
	group.Post("/", middleware.RBACMiddleware("achievement:create"), achievementService.CreateAchievement)
	group.Post("/:id/attachments", middleware.RBACMiddleware("achievement:update"), achievementService.UploadAttachment)
	group.Put("/:id", middleware.RBACMiddleware("achievement:update"), achievementService.UpdateAchievement)
	group.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), achievementService.DeleteAchievement)
	group.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), achievementService.SubmitAchievement)
	group.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), achievementService.VerifyAchievement)
	group.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), achievementService.RejectAchievement)
	group.Get("/advisee/list", middleware.RBACMiddleware("achievement:read"), achievementService.GetAdviseeAchievements)
}

func SetupLecturerRoutes(app *fiber.App, lecturerService service.LecturerService) {
	group := app.Group("/api/v1/lecturers", middleware.AuthMiddleware())

	group.Get("/", middleware.RBACMiddleware("lecturer:read"), lecturerService.GetAllLecturers)
	group.Get("/:id", middleware.RBACMiddleware("lecturer:read"), lecturerService.GetLecturerByID)
	group.Get("/:id/advisees", middleware.RBACMiddleware("lecturer:read"), lecturerService.GetAdvisees)
	group.Post("/", middleware.RBACMiddleware("lecturer:create"), lecturerService.CreateLecturer)
	group.Put("/:id", middleware.RBACMiddleware("lecturer:update"), lecturerService.UpdateLecturer)
	group.Delete("/:id", middleware.RBACMiddleware("lecturer:delete"), lecturerService.DeleteLecturer)
}

func SetupRoleRoutes(app *fiber.App, roleService service.RoleService) {
	group := app.Group("/api/v1/roles", middleware.AuthMiddleware())

	group.Get("/", middleware.RBACMiddleware("role:read"), roleService.GetAllRoles)
	group.Get("/:id", middleware.RBACMiddleware("role:read"), roleService.GetRoleByID)
	group.Get("/name/:name", middleware.RBACMiddleware("role:read"), roleService.GetRoleByName)
	group.Post("/", middleware.RBACMiddleware("role:create"), roleService.CreateRole)
	group.Put("/:id", middleware.RBACMiddleware("role:update"), roleService.UpdateRole)
	group.Post("/assign-permission", middleware.RBACMiddleware("role:assign-permission"), roleService.AssignPermission)
	group.Post("/remove-permission", middleware.RBACMiddleware("role:remove-permission"), roleService.RemovePermission)
}

func SetupStudentRoutes(app *fiber.App, studentService service.StudentService) {
	group := app.Group("/api/v1/students", middleware.AuthMiddleware())

	group.Get("/", middleware.RBACMiddleware("student:read"), studentService.GetAllStudents)
	group.Get("/:id", middleware.RBACMiddleware("student:read"), studentService.GetStudentByID)
	group.Get("/:id/achievements", middleware.RBACMiddleware("student:read"), studentService.GetStudentAchievements)
	group.Get("/user/:user_id", middleware.RBACMiddleware("student:read"), studentService.GetStudentByUserID)
	group.Get("/advisor/:advisor_id", middleware.RBACMiddleware("student:read"), studentService.GetStudentsByAdvisor)
	group.Post("/", middleware.RBACMiddleware("student:create"), studentService.CreateStudent)
	group.Put("/:id", middleware.RBACMiddleware("student:update"), studentService.UpdateStudent)
	group.Put("/:id/advisor", middleware.RBACMiddleware("student:update"), studentService.UpdateAdvisor)
	group.Delete("/:id", middleware.RBACMiddleware("student:delete"), studentService.DeleteStudent)
}

func SetupUserRoutes(app *fiber.App, userService service.UserService) {
	group := app.Group("/api/v1/users", middleware.AuthMiddleware())

	group.Get("/", middleware.RBACMiddleware("user:read"), userService.GetAllUsers)
	group.Get("/:id", middleware.RBACMiddleware("user:read"), userService.GetUserByID)
	group.Post("/", middleware.RBACMiddleware("user:create"), userService.CreateUser)
	group.Put("/:id", middleware.RBACMiddleware("user:update"), userService.UpdateUser)
	group.Put("/:id/role", middleware.RBACMiddleware("user:update"), userService.AssignRole)
	group.Delete("/:id", middleware.RBACMiddleware("user:delete"), userService.DeleteUser)
}

func SetupReportRoutes(app *fiber.App, reportService service.ReportService) {
	group := app.Group("/api/v1/reports", middleware.AuthMiddleware())

	group.Get("/statistics", middleware.RBACMiddleware("report:read"), reportService.GetStatistics)
	group.Get("/statistics/period", middleware.RBACMiddleware("report:read"), reportService.GetStatisticsByPeriod)
	group.Get("/statistics/type", middleware.RBACMiddleware("report:read"), reportService.GetStatisticsByType)
	group.Get("/top-students", middleware.RBACMiddleware("report:read"), reportService.GetTopStudents)
	group.Get("/student/:id", middleware.RBACMiddleware("report:read"), reportService.GetStudentReport)
}

func SetupPermissionRoutes(app *fiber.App, permissionService service.PermissionService) {
	group := app.Group("/api/v1/permissions", middleware.AuthMiddleware())

	group.Get("/", middleware.RBACMiddleware("permission:read"), permissionService.GetAllPermissions)
	group.Post("/", middleware.RBACMiddleware("permission:create"), permissionService.CreatePermission)
	group.Get("/role/:roleId", middleware.RBACMiddleware("permission:read"), permissionService.GetPermissionsByRoleID)
}

func SetupAuthRoutes(app *fiber.App, authService service.AuthService) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.RefreshToken)

	protected := auth.Group("", middleware.AuthMiddleware())
	protected.Get("/profile", authService.GetProfile)
	protected.Post("/logout", authService.Logout)
}
