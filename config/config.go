package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config menyimpan semua konfigurasi aplikasi dari environment variables
type Config struct {
	Database DatabaseConfig
	MongoDB  MongoDBConfig // Add MongoDB configuration
	Server   ServerConfig
	JWT      JWTConfig
}

// DatabaseConfig menyimpan konfigurasi PostgreSQL database
type DatabaseConfig struct {
	Host     string // DB_HOST - hostname database (default: localhost)
	Port     string // DB_PORT - port database (default: 5432)
	User     string // DB_USER - username database (default: postgres)
	Password string // DB_PASSWORD - password database (default: postgres)
	DBName   string // DB_NAME - nama database (default: uas_be_db)
	SSLMode  string // DB_SSLMODE - SSL mode connection (default: disable)
}

// MongoDBConfig menyimpan konfigurasi MongoDB untuk data prestasi
type MongoDBConfig struct {
	URI        string // MONGO_URI - MongoDB connection URI
	Database   string // MONGO_DB_NAME - nama database MongoDB (default: uas_achievements_db)
	Collection string // Collection name for achievements (default: achievements)
}

// ServerConfig menyimpan konfigurasi server Fiber
type ServerConfig struct {
	Host string // SERVER_HOST - host server (default: 0.0.0.0)
	Port string // SERVER_PORT - port server (default: 8080)
}

// JWTConfig menyimpan konfigurasi JWT
type JWTConfig struct {
	Secret string // JWT_SECRET - secret key untuk JWT (default: mysecretkey)
}

// LoadConfig memuat konfigurasi dari environment variables dengan default values
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "postgres"),
			DBName:   GetEnv("DB_NAME", "uas_be_db"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		MongoDB: MongoDBConfig{
			URI:        GetEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database:   GetEnv("MONGO_DB_NAME", "uas_achievements_db"),
			Collection: "achievements",
		},
		Server: ServerConfig{
			Host: GetEnv("SERVER_HOST", "0.0.0.0"),
			Port: GetEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret: GetEnv("JWT_SECRET", "mysecretkey"),
		},
	}
}

// GetDSN mengembalikan connection string (Data Source Name) untuk database PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}

// GetEnv mengambil environment variable dengan fallback ke default value jika tidak ada
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt mengambil environment variable sebagai integer dengan default value
func getEnvAsInt(name string, defaultVal int) int {
	valStr := GetEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
