package model

// Permission merepresentasikan permission yang bisa dilakukan dalam sistem
type Permission struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`               // Nama permission (contoh: achievement:create)
	Resource    string `db:"resource" json:"resource"`       // Resource yang diproteksi (achievement, user)
	Action      string `db:"action" json:"action"`           // Aksi yang diproteksi (create, read, update, delete, verify)
	Description string `db:"description" json:"description"` // Deskripsi permission
}
