package main

import (
	"context"
	"database/sql"
	"log"

	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/config"
	"uas_be/database"
	_ "uas_be/docs"
	"uas_be/route"
	"uas_be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title UAS Backend API - Sistem Pelaporan Prestasi Mahasiswa
// @version 1.0
// @description API untuk mengelola prestasi mahasiswa dengan fitur autentikasi, RBAC, dan pelaporan
// @description Sistem ini mengelola data prestasi mahasiswa, dosen wali, dan admin dengan kontrol akses berbasis peran
// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found or failed to load; continuing with environment variables or defaults")
	}

	cfg := config.LoadConfig()

	log.Printf("📝 MongoDB Config - URI: %s, Database: %s", cfg.MongoDB.URI, cfg.MongoDB.Database)

	// Initialize JWT secret
	utils.InitJWT(cfg.JWT.Secret)

	db := database.InitPostgres(cfg)
	if err := database.InitSchema(db); err != nil {
		log.Println("warning: failed to init schema:", err)
	}

	// Initialize MongoDB for achievements data
	_, mongoDB, err := database.InitMongoDB(cfg)
	if err != nil {
		log.Fatal("❌ MongoDB connection failed:", err)
	}

	collections, err := mongoDB.ListCollectionNames(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Println("⚠️  Warning: Failed to list collections:", err)
	} else {
		log.Printf("📦 Available collections in '%s': %v", cfg.MongoDB.Database, collections)
	}

	database.SetDB(db)

	// Seed default admin user
	seedDefaultAdmin(db)

	// ===== FIBER APP =====
	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	route.RegisterRoutes(app)

	log.Println("🚀 Server running at http://localhost:3000")
	log.Println("📚 Swagger documentation available at http://localhost:3000/swagger/index.html")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func seedDefaultAdmin(db *sql.DB) {
	// Check if admin user already exists
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	adminRole, err := roleRepo.GetRoleByName("Admin")
	if err != nil {
		log.Println("Warning: failed to get admin role:", err)
		return
	}
	if adminRole == nil {
		log.Println("Warning: admin role not found")
		return
	}

	existingAdmin, err := userRepo.GetUserByUsername("admin")
	if err != nil {
		log.Println("Warning: failed to check existing admin:", err)
		return
	}

	// Create default admin
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Println("Warning: failed to hash admin password:", err)
		return
	}

	if existingAdmin != nil {
		// Update password if exists
		existingAdmin.PasswordHash = hashedPassword
		err = userRepo.UpdateUser(existingAdmin)
		if err != nil {
			log.Println("Warning: failed to update admin password:", err)
			return
		}
		log.Println("✅ Admin password updated: username=admin, password=admin123")
	} else {
		adminUser := &model.User{
			Username:     "admin",
			Email:        "admin@uas.com",
			PasswordHash: hashedPassword,
			FullName:     "Administrator",
			RoleID:       adminRole.ID,
			IsActive:     true,
		}

		err = userRepo.CreateUser(adminUser)
		if err != nil {
			log.Println("Warning: failed to create default admin:", err)
			return
		}

		log.Println("✅ Default admin user created: username=admin, password=admin123")
	}
}
