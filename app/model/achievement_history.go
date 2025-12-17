package model

import "time"

// AchievementHistory menyimpan riwayat perubahan status achievement
type AchievementHistory struct {
	ID            string    `db:"id" json:"id"`
	AchievementID string    `db:"achievement_id" json:"achievement_id"`
	OldStatus     string    `db:"previous_status" json:"previous_status"`
	NewStatus     string    `db:"new_status" json:"new_status"`
	ChangedBy     string    `db:"changed_by" json:"changed_by"`
	ChangedByName string    `db:"changed_by_name" json:"changed_by_name"`
	Note          *string   `db:"notes" json:"notes"`
	CreatedAt     time.Time `db:"changed_at" json:"changed_at"`
}
