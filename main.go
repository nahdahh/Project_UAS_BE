package main

import (
	"fmt"
	"log"

	"uas_be/app/repository"
	"uas_be/app/service"
	"uas_be/config"
	"uas_be/database"
	"uas_be/route"
	"uas_be/utils"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// ===== STEP 1: LOAD CONFIGURATION =====
	// Muat konfigurasi dari environment variables
	cfg := config.LoadConfig()
	log.Println("✅ Konfigurasi berhasil dimuat")

	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key-change-in-production")
	utils.InitJWT(jwtSecret)
	log.Println("✅ JWT secret berhasil diinisialisasi")

	//  INITIALIZE DATABASE =====
	// Inisialisasi koneksi PostgreSQL
	pgDB := database.InitPostgres(cfg)
	defer pgDB.Close()

	// Inisialisasi schema dan tabel di database
	if err := database.InitSchema(pgDB); err != nil {
		log.Fatalf("❌ Gagal inisialisasi schema: %v\n", err)
	}

	//  INITIALIZE REPOSITORIES =====
	// Repository adalah layer untuk akses data dari database
	userRepo := repository.NewUserRepository(pgDB)
	roleRepo := repository.NewRoleRepository(pgDB)
	permissionRepo := repository.NewPermissionRepository(pgDB)
	studentRepo := repository.NewStudentRepository(pgDB)
	lecturerRepo := repository.NewLecturerRepository(pgDB)
	achievementRepo := repository.NewAchievementRepository(pgDB)
	log.Println("✅ Repositories berhasil diinisialisasi")

	// =====  INITIALIZE SERVICES =====

	authService := service.NewAuthService(userRepo, permissionRepo)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	studentService := service.NewStudentService(studentRepo)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo)
	achievementService := service.NewAchievementService(achievementRepo, studentRepo)
	log.Println("✅ Services berhasil diinisialisasi")

	// ===== INITIALIZE FIBER APP =====
	// Setup Fiber app dengan konfigurasi dasar
	app := fiber.New(fiber.Config{
		AppName: "UAS_BE - Sistem Pelaporan Prestasi Mahasiswa v1.0",
	})

	// =====  REGISTER ROUTES =====
	// Daftarkan semua API routes
	route.RegisterRoutes(
		app,
		authService,
		userService,
		studentService,
		lecturerService,
		achievementService,
		roleService,
	)
	log.Println("✅ Routes berhasil didaftarkan")

	// =====  START SERVER =====
	// Jalankan server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 Server UAS_BE berjalan di http://localhost:%s\n", port)
	fmt.Println("📚 API Documentation: http://localhost:" + port + "/health")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
