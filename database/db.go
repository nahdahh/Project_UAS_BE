package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"uas_be/config"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global db variables
var db *sql.DB
var mongoClient *mongo.Client // Add MongoDB client
var mongoDB *mongo.Database   // Add MongoDB database reference

// SetDB sets the global database connection
func SetDB(database *sql.DB) {
	db = database
}

// GetDB returns the global database connection
func GetDB() *sql.DB {
	return db
}

// SetMongoDB sets the global MongoDB connection
func SetMongoDB(client *mongo.Client, database *mongo.Database) {
	mongoClient = client
	mongoDB = database
}

// GetMongoDB returns the global MongoDB database
func GetMongoDB() *mongo.Database {
	return mongoDB
}

// GetMongoClient returns the global MongoDB client
func GetMongoClient() *mongo.Client {
	return mongoClient
}

// InitPostgres membuka koneksi PostgreSQL
func InitPostgres(cfg *config.Config) *sql.DB {
	dsn := cfg.Database.GetDSN()
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Gagal membuka koneksi PostgreSQL:", err)
	}

	// Test koneksi database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbConn.PingContext(ctx); err != nil {
		log.Fatal("❌ Gagal ping database:", err)
	}

	log.Println("✅ Berhasil terhubung ke PostgreSQL")
	SetDB(dbConn) // Set the global database connection
	return dbConn
}

// InitMongoDB membuka koneksi MongoDB untuk menyimpan data prestasi
func InitMongoDB(cfg *config.Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Test koneksi MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	database := client.Database(cfg.MongoDB.Database)
	SetMongoDB(client, database)

	if err := InitMongoCollections(database, cfg.MongoDB.Collection); err != nil {
		log.Println("⚠️  Warning: Failed to initialize MongoDB collections:", err)
	}

	log.Printf("✅ Berhasil terhubung ke MongoDB - Database: %s", cfg.MongoDB.Database)
	return client, database, nil
}

// InitMongoCollections membuat collection dan indexes di MongoDB
func InitMongoCollections(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get atau create collection
	collection := db.Collection(collectionName)

	// Create indexes untuk performa query
	indexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"studentId": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"achievementType": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"createdAt": -1,
			},
		},
		{
			Keys: map[string]interface{}{
				"tags": 1,
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	log.Printf("✅ MongoDB collection '%s' siap digunakan dengan indexes", collectionName)
	return nil
}

// InitSchema membuat schema dan tabel awal di database
func InitSchema(db *sql.DB) error {
	schema := `
	-- Tabel roles: menyimpan daftar role (Admin, Mahasiswa, Dosen Wali)
	CREATE TABLE IF NOT EXISTS roles (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(50) UNIQUE NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Tabel permissions: menyimpan daftar permission yang bisa dilakukan
	CREATE TABLE IF NOT EXISTS permissions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) UNIQUE NOT NULL,
		resource VARCHAR(50) NOT NULL,
		action VARCHAR(50) NOT NULL,
		description TEXT
	);

	-- Tabel role_permissions: menghubungkan role dengan permissions (many-to-many)
	CREATE TABLE IF NOT EXISTS role_permissions (
		role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
		permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
		PRIMARY KEY (role_id, permission_id)
	);

	-- Tabel users: menyimpan data pengguna sistem
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		full_name VARCHAR(100) NOT NULL,
		role_id UUID NOT NULL REFERENCES roles(id),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Tabel lecturers: menyimpan data dosen
	CREATE TABLE IF NOT EXISTS lecturers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
		lecturer_id VARCHAR(20) UNIQUE NOT NULL,
		department VARCHAR(100),
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Tabel students: menyimpan data mahasiswa
	CREATE TABLE IF NOT EXISTS students (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
		student_id VARCHAR(20) UNIQUE NOT NULL,
		program_study VARCHAR(100),
		academic_year VARCHAR(10),
		advisor_id UUID REFERENCES lecturers(id),
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Update achievement_references table - this is now the primary reference table
	-- Tabel achievement_references: menyimpan referensi prestasi (pointer ke MongoDB)
	CREATE TABLE IF NOT EXISTS achievement_references (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		mongo_achievement_id VARCHAR(24) NOT NULL,
		achievement_title VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'draft',
		submitted_at TIMESTAMP,
		verified_at TIMESTAMP,
		verified_by UUID REFERENCES users(id),
		rejection_note TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Tabel achievement_history: menyimpan riwayat perubahan prestasi
	CREATE TABLE IF NOT EXISTS achievement_history (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		achievement_id UUID NOT NULL REFERENCES achievement_references(id),
		changed_by UUID NOT NULL REFERENCES users(id),
		action VARCHAR(50) NOT NULL,
		previous_status VARCHAR(20),
		new_status VARCHAR(20),
		notes TEXT,
		changed_at TIMESTAMP DEFAULT NOW()
	);

	-- Tabel achievement_attachments: menyimpan file lampiran prestasi
	CREATE TABLE IF NOT EXISTS achievement_attachments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		achievement_id UUID NOT NULL REFERENCES achievement_references(id),
		file_name VARCHAR(255) NOT NULL,
		file_path VARCHAR(500) NOT NULL,
		file_type VARCHAR(50),
		file_size BIGINT,
		uploaded_by UUID NOT NULL REFERENCES users(id),
		uploaded_at TIMESTAMP DEFAULT NOW()
	);

	-- Masukkan data awal untuk roles
	INSERT INTO roles (name, description) VALUES 
		('Admin', 'Administrator sistem dengan akses penuh'),
		('Mahasiswa', 'Pengguna mahasiswa'),
		('Dosen Wali', 'Dosen pembimbing akademik')
	ON CONFLICT (name) DO NOTHING;

	-- Masukkan data awal untuk permissions
	INSERT INTO permissions (name, resource, action, description) VALUES
		('achievement:create', 'achievement', 'create', 'Membuat prestasi baru'),
		('achievement:read', 'achievement', 'read', 'Membaca prestasi'),
		('achievement:update', 'achievement', 'update', 'Mengubah prestasi'),
		('achievement:delete', 'achievement', 'delete', 'Menghapus prestasi'),
		('achievement:verify', 'achievement', 'verify', 'Memverifikasi prestasi'),
		('achievement:submit', 'achievement', 'submit', 'Mengajukan prestasi'),
		('user:create', 'user', 'create', 'Membuat pengguna baru'),
		('user:read', 'user', 'read', 'Membaca data pengguna'),
		('user:update', 'user', 'update', 'Mengubah data pengguna'),
		('user:delete', 'user', 'delete', 'Menghapus pengguna'),
		('lecturer:create', 'lecturer', 'create', 'Membuat data dosen'),
		('lecturer:read', 'lecturer', 'read', 'Membaca data dosen'),
		('lecturer:update', 'lecturer', 'update', 'Mengubah data dosen'),
		('lecturer:delete', 'lecturer', 'delete', 'Menghapus data dosen'),
		('student:create', 'student', 'create', 'Membuat data mahasiswa'),
		('student:read', 'student', 'read', 'Membaca data mahasiswa'),
		('student:update', 'student', 'update', 'Mengubah data mahasiswa'),
		('student:delete', 'student', 'delete', 'Menghapus data mahasiswa'),
		('role:create', 'role', 'create', 'Membuat role baru'),
		('role:read', 'role', 'read', 'Membaca data role'),
		('role:update', 'role', 'update', 'Mengubah data role'),
		('role:assign-permission', 'role', 'assign-permission', 'Menetapkan permission ke role'),
		('role:remove-permission', 'role', 'remove-permission', 'Menghapus permission dari role')
	ON CONFLICT (name) DO NOTHING;

	-- Assign permissions ke role Admin (semua permission)
	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'Admin'
	ON CONFLICT DO NOTHING;

	-- Assign permissions ke role Mahasiswa
	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id FROM roles r, permissions p 
	WHERE r.name = 'Mahasiswa' AND p.name IN ('achievement:create', 'achievement:read', 'achievement:update', 'achievement:delete')
	ON CONFLICT DO NOTHING;

	-- Assign permissions ke role Dosen Wali
	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id FROM roles r, permissions p 
	WHERE r.name = 'Dosen Wali' AND p.name IN ('achievement:read', 'achievement:verify')
	ON CONFLICT DO NOTHING;
	`

	_, err := db.Exec(schema)
	if err != nil {
		log.Println("❌ Gagal membuat schema:", err)
		return err
	}

	log.Println("✅ Schema berhasil dibuat")
	return nil
}
