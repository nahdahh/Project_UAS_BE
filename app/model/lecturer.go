package model

import "time"

// Lecturer merepresentasikan data dosen
type Lecturer struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`         // Foreign key ke tabel users
	LecturerID string    `db:"lecturer_id" json:"lecturer_id"` // NIDN atau identitas dosen
	Department string    `db:"department" json:"department"`   // Departemen/Jurusan
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// LecturerWithUser menampilkan data lecturer dengan data user
type LecturerWithUser struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	LecturerID string `json:"lecturer_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	CreatedAt  string `json:"created_at"`
}
