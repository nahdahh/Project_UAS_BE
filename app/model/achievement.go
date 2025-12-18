package model

import "time"

// Achievement status constants
const (
	AchievementStatusDraft     = "draft"
	AchievementStatusSubmitted = "submitted"
	AchievementStatusVerified  = "verified"
	AchievementStatusRejected  = "rejected"
	AchievementStatusDeleted   = "deleted" // Added "deleted" status for soft delete
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

// Achievement menyimpan prestasi mahasiswa di MongoDB
type Achievement struct {
	ID              string                 `bson:"_id,omitempty" json:"id"`
	StudentID       string                 `bson:"student_id" json:"student_id"`
	AchievementType string                 `bson:"achievement_type" json:"achievement_type"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	Details         map[string]interface{} `bson:"details" json:"details"` // Dynamic fields
	Tags            []string               `bson:"tags" json:"tags"`
	Points          int                    `bson:"points" json:"points"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}

// AchievementReference menyimpan referensi prestasi di PostgreSQL (pointer ke MongoDB)
type AchievementReference struct {
	ID                 string     `db:"id" json:"id"`
	StudentID          string     `db:"student_id" json:"student_id"`
	MongoAchievementID string     `db:"mongo_achievement_id" json:"mongo_achievement_id"` // Pointer ke MongoDB ObjectId
	AchievementTitle   string     `db:"achievement_title" json:"achievement_title"`
	Status             string     `db:"status" json:"status"`
	SubmittedAt        *time.Time `db:"submitted_at" json:"submitted_at"`
	VerifiedAt         *time.Time `db:"verified_at" json:"verified_at"`
	VerifiedBy         *string    `db:"verified_by" json:"verified_by"`
	RejectionNote      *string    `db:"rejection_note" json:"rejection_note"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at"` // Added deleted_at field
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

// AchievementWithReference combines MongoDB data with PostgreSQL reference
type AchievementWithReference struct {
	Achievement
	Status        string     `json:"status"`
	SubmittedAt   *time.Time `json:"submitted_at"`
	VerifiedAt    *time.Time `json:"verified_at"`
	VerifiedBy    *string    `json:"verified_by"`
	RejectionNote *string    `json:"rejection_note"`
}

// CreateAchievementRequest adalah request untuk membuat prestasi baru
type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type"` // Tipe prestasi: academic, competition, organization, publication, certification, other
	Title           string                 `json:"title"`            // Judul prestasi
	Description     string                 `json:"description"`      // Deskripsi prestasi
	Details         map[string]interface{} `json:"details"`          // Detail dinamis berdasarkan tipe
	Tags            []string               `json:"tags"`             // Tag/kategori prestasi
	// Points field removed - points will be assigned by lecturer during verification
}

// UpdateAchievementRequest adalah request untuk update prestasi
type UpdateAchievementRequest struct {
	AchievementType *string                 `json:"achievement_type"`
	Title           *string                 `json:"title"`
	Description     *string                 `json:"description"`
	Details         *map[string]interface{} `json:"details"`
	Tags            *[]string               `json:"tags"`
	// Points field removed - points cannot be updated by students
}

// VerifyAchievementRequest adalah request untuk verifikasi prestasi dengan poin
type VerifyAchievementRequest struct {
	Points int `json:"points" validate:"required,min=0"` // Poin yang diberikan dosen saat verifikasi
}
