package model

import "time"

// Role merepresentasikan role pengguna (Admin, Mahasiswa, Dosen Wali)
type Role struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`               // Nama role (Admin, Mahasiswa, Dosen Wali)
	Description string    `db:"description" json:"description"` // Deskripsi role
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
