package model

import "time"

// Achievement status constants
const (
	AchievementStatusDraft     = "draft"
	AchievementStatusSubmitted = "submitted"
	AchievementStatusVerified  = "verified"
	AchievementStatusRejected  = "rejected"
)

// AchievementReference menyimpan referensi prestasi mahasiswa di PostgreSQL
type AchievementReference struct {
	ID               string     `db:"id" json:"id"`
	StudentID        string     `db:"student_id" json:"student_id"`               // Foreign key ke students
	AchievementTitle string     `db:"achievement_title" json:"achievement_title"` // Judul prestasi
	Status           string     `db:"status" json:"status"`                       // Status: draft, submitted, verified, rejected
	SubmittedAt      *time.Time `db:"submitted_at" json:"submitted_at"`           // Waktu submit untuk verifikasi
	VerifiedAt       *time.Time `db:"verified_at" json:"verified_at"`             // Waktu verifikasi
	VerifiedBy       *string    `db:"verified_by" json:"verified_by"`             // ID dosen yang memverifikasi
	RejectionNote    *string    `db:"rejection_note" json:"rejection_note"`       // Catatan penolakan
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// CreateAchievementRequest adalah request untuk membuat prestasi baru
type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type"` // Tipe prestasi: academic, competition, organization, publication, certification, other
	Title           string                 `json:"title"`            // Judul prestasi
	Description     string                 `json:"description"`      // Deskripsi prestasi
	Details         map[string]interface{} `json:"details"`          // Detail dinamis berdasarkan tipe
	Tags            []string               `json:"tags"`             // Tag/kategori prestasi
	Points          int                    `json:"points"`           // Poin prestasi
}

// UpdateAchievementRequest adalah request untuk update prestasi
type UpdateAchievementRequest struct {
	AchievementType *string                 `json:"achievement_type"`
	Title           *string                 `json:"title"`
	Description     *string                 `json:"description"`
	Details         *map[string]interface{} `json:"details"`
	Tags            *[]string               `json:"tags"`
	Points          *int                    `json:"points"`
}

// AchievementWithDetails menampilkan achievement dengan detail lengkap
type AchievementWithDetails struct {
	ID              string                 `json:"id"`
	StudentID       string                 `json:"student_id"`
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
	Status          string                 `json:"status"`
	SubmittedAt     *time.Time             `json:"submitted_at"`
	VerifiedAt      *time.Time             `json:"verified_at"`
	VerifiedBy      *string                `json:"verified_by"`
	RejectionNote   *string                `json:"rejection_note"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
