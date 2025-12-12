package main

import (
	"log"

	"uas_be/config"
	"uas_be/database"
	"uas_be/route"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found or failed to load; continuing with environment variables or defaults")
	}

	cfg := config.LoadConfig()

	db := database.InitPostgres(cfg)
	if err := database.InitSchema(db); err != nil {
		log.Println("warning: failed to init schema:", err)
	}

	database.SetDB(db)

	// ===== FIBER APP =====
	app := fiber.New()

	route.RegisterRoutes(app)

	log.Println("🚀 Server running at http://localhost:3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
