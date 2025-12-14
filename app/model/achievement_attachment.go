package model

import "time"

// AchievementAttachment menyimpan attachment/file untuk achievement
type AchievementAttachment struct {
	ID            string    `db:"id" json:"id"`
	AchievementID string    `db:"achievement_id" json:"achievement_id"`
	FileName      string    `db:"file_name" json:"file_name"`
	FileURL       string    `db:"file_url" json:"file_url"`
	FileSize      int64     `db:"file_size" json:"file_size"`
	FileType      string    `db:"file_type" json:"file_type"`
	UploadedBy    string    `db:"uploaded_by" json:"uploaded_by"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
