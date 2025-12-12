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
	studentService := service.NewStudentService(studentRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	userService := service.NewUserService(userRepo)
	
	SetupAuthRoutes(app, authService)
	SetupAchievementRoutes(app, achievementService)
	SetupLecturerRoutes(app, lecturerService)
	SetupStudentRoutes(app, studentService)
	SetupRoleRoutes(app, roleService)
	SetupUserRoutes(app, userService)

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

	// List achievements (filtered by role)
	group.Get("/", middleware.RBACMiddleware("achievement:read"), achievementService.GetAllAchievements)

	// Detail achievement
	group.Get("/:id", middleware.RBACMiddleware("achievement:read"), achievementService.GetAchievementDetail)

	// Create achievement
	group.Post("/", middleware.RBACMiddleware("achievement:create"), achievementService.CreateAchievement)

	// Update achievement
	group.Put("/:id", middleware.RBACMiddleware("achievement:update"), achievementService.UpdateAchievement)

	// Delete achievement
	group.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), achievementService.DeleteAchievement)

	// Submit for verification
	group.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), achievementService.SubmitAchievement)

	// Verify achievement
	group.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), achievementService.VerifyAchievement)

	// Reject achievement
	group.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), achievementService.RejectAchievement)

	// Get advisee achievements
	group.Get("/advisee/list", middleware.RBACMiddleware("achievement:read"), achievementService.GetAdviseeAchievements)
}

func SetupLecturerRoutes(app *fiber.App, lecturerService service.LecturerService) {
	group := app.Group("/api/v1/lecturers", middleware.AuthMiddleware())

	// List lecturers
	group.Get("/", middleware.RBACMiddleware("lecturer:read"), lecturerService.GetAllLecturers)

	// Get lecturer detail
	group.Get("/:id", middleware.RBACMiddleware("lecturer:read"), lecturerService.GetLecturerByID)

	// Get advisees
	group.Get("/:id/advisees", middleware.RBACMiddleware("lecturer:read"), lecturerService.GetAdvisees)

	// Create lecturer
	group.Post("/", middleware.RBACMiddleware("lecturer:create"), lecturerService.CreateLecturer)

	// Update lecturer
	group.Put("/:id", middleware.RBACMiddleware("lecturer:update"), lecturerService.UpdateLecturer)

	// Delete lecturer
	group.Delete("/:id", middleware.RBACMiddleware("lecturer:delete"), lecturerService.DeleteLecturer)
}

func SetupRoleRoutes(app *fiber.App, roleService service.RoleService) {
	group := app.Group("/api/v1/roles", middleware.AuthMiddleware())

	// List roles
	group.Get("/", middleware.RBACMiddleware("role:read"), roleService.GetAllRoles)

	// Get role by ID
	group.Get("/:id", middleware.RBACMiddleware("role:read"), roleService.GetRoleByID)

	// Get role by name
	group.Get("/name/:name", middleware.RBACMiddleware("role:read"), roleService.GetRoleByName)

	// Create role
	group.Post("/", middleware.RBACMiddleware("role:create"), roleService.CreateRole)

	// Update role
	group.Put("/:id", middleware.RBACMiddleware("role:update"), roleService.UpdateRole)

	// Assign permission
	group.Post("/assign-permission", middleware.RBACMiddleware("role:assign-permission"), roleService.AssignPermission)

	// Remove permission
	group.Post("/remove-permission", middleware.RBACMiddleware("role:remove-permission"), roleService.RemovePermission)
}

func SetupStudentRoutes(app *fiber.App, studentService service.StudentService) {
	group := app.Group("/api/v1/students", middleware.AuthMiddleware())

	// List students
	group.Get("/", middleware.RBACMiddleware("student:read"), studentService.GetAllStudents)

	// Get student by ID
	group.Get("/:id", middleware.RBACMiddleware("student:read"), studentService.GetStudentByID)

	// Get student by user ID
	group.Get("/user/:user_id", middleware.RBACMiddleware("student:read"), studentService.GetStudentByUserID)

	// Get students by advisor
	group.Get("/advisor/:advisor_id", middleware.RBACMiddleware("student:read"), studentService.GetStudentsByAdvisor)

	// Create student
	group.Post("/", middleware.RBACMiddleware("student:create"), studentService.CreateStudent)

	// Update student
	group.Put("/:id", middleware.RBACMiddleware("student:update"), studentService.UpdateStudent)

	// Delete student
	group.Delete("/:id", middleware.RBACMiddleware("student:delete"), studentService.DeleteStudent)
}

func SetupUserRoutes(app *fiber.App, userService service.UserService) {
	group := app.Group("/api/v1/users", middleware.AuthMiddleware())

	// List users
	group.Get("/", middleware.RBACMiddleware("user:read"), userService.GetAllUsers)

	// Get user by ID
	group.Get("/:id", middleware.RBACMiddleware("user:read"), userService.GetUserByID)

	// Create user
	group.Post("/", middleware.RBACMiddleware("user:create"), userService.CreateUser)

	// Update user
	group.Put("/:id", middleware.RBACMiddleware("user:update"), userService.UpdateUser)

	// Delete user
	group.Delete("/:id", middleware.RBACMiddleware("user:delete"), userService.DeleteUser)
}

func SetupAuthRoutes(app *fiber.App, authService service.AuthService) {
	group := app.Group("/api/v1/auth", middleware.AuthMiddleware())

	// Login route
	group.Post("/login", authService.Login)

	// Logout route
	group.Post("/logout", authService.Logout)
}
