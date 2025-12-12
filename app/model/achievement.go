package model

import "time"

// Achievement status constants
const (
	AchievementStatusDraft     = "draft"
	AchievementStatusSubmitted = "submitted"
	AchievementStatusVerified  = "verified"
	AchievementStatusRejected  = "rejected"
)

// Achievement type constants
const (
	AchievementTypeAcademic      = "academic"
	AchievementTypeCompetition   = "competition"
	AchievementTypeOrganization  = "organization"
	AchievementTypePublication   = "publication"
	AchievementTypeCertification = "certification"
	AchievementTypeOther         = "other"
)

// Achievement menyimpan prestasi mahasiswa
type Achievement struct {
	ID              string                 `db:"id" json:"id"`
	StudentID       string                 `db:"student_id" json:"student_id"`
	AchievementType string                 `db:"achievement_type" json:"achievement_type"` // academic, competition, organization, publication, certification, other
	Title           string                 `db:"title" json:"title"`
	Description     string                 `db:"description" json:"description"`
	Details         map[string]interface{} `db:"details" json:"details"`   // JSONB - detail dinamis berdasarkan tipe
	Tags            []string               `db:"tags" json:"tags"`         // Array - tag/kategori prestasi
	Points          int                    `db:"points" json:"points"`     // Poin prestasi
	Status          string                 `db:"status" json:"status"`     // Status: draft, submitted, verified, rejected
	SubmittedAt     *time.Time             `db:"submitted_at" json:"submitted_at"`
	VerifiedAt      *time.Time             `db:"verified_at" json:"verified_at"`
	VerifiedBy      *string                `db:"verified_by" json:"verified_by"`
	RejectionNote   *string                `db:"rejection_note" json:"rejection_note"`
	CreatedAt       time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `db:"updated_at" json:"updated_at"`
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
