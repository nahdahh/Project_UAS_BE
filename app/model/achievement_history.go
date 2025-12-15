package model

import "time"

// AchievementHistory menyimpan riwayat perubahan status achievement
type AchievementHistory struct {
	ID            string    `db:"id" json:"id"`
	AchievementID string    `db:"achievement_id" json:"achievement_id"`
	OldStatus     string    `db:"old_status" json:"old_status"`
	NewStatus     string    `db:"new_status" json:"new_status"`
	ChangedBy     string    `db:"changed_by" json:"changed_by"`
	ChangedByName string    `db:"changed_by_name" json:"changed_by_name"`
	Note          *string   `db:"note" json:"note"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
