package model

import "time"

// User merepresentasikan pengguna dalam sistem
type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`           // Username unik pengguna
	Email        string    `db:"email" json:"email"`                 // Email unik pengguna
	PasswordHash string    `db:"password_hash" json:"password_hash"` // Hash password (tidak dikembalikan ke client)
	FullName     string    `db:"full_name" json:"full_name"`         // Nama lengkap pengguna
	RoleID       string    `db:"role_id" json:"role_id"`             // Foreign key ke tabel roles
	IsActive     bool      `db:"is_active" json:"is_active"`         // Status aktif pengguna
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// UserWithRole menampilkan user dengan role name dan permissions
type UserWithRole struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	RoleName    string   `json:"role"`
	IsActive    bool     `json:"is_active"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
}

// CreateUserRequest adalah request untuk membuat user baru
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	RoleID   string `json:"role_id"`
}

// LoginRequest adalah request untuk login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse adalah response setelah login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserWithRole `json:"user"`
}

// UserProfile adalah profile user yang dikembalikan ke client
type UserProfile struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	RoleName    string   `json:"role"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
}
