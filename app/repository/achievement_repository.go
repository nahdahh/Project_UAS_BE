package repository

import (
	"database/sql"
	"encoding/json"
	"uas_be/app/model"

	"github.com/lib/pq"
)

// AchievementRepository adalah interface untuk akses data achievement dari database
type AchievementRepository interface {
	CreateAchievement(achievement *model.Achievement) error
	GetAchievementByID(id string) (*model.Achievement, error)
	GetAchievementsByStudentID(studentID string) ([]*model.Achievement, error)
	GetAchievementsByStatus(status string) ([]*model.Achievement, error)
	GetAllAchievements(page, pageSize int) ([]*model.Achievement, int, error)
	UpdateAchievement(achievement *model.Achievement) error
	SubmitAchievementForVerification(id string) error
	VerifyAchievement(id string, verifiedBy string) error
	RejectAchievement(id string, verifiedBy string, rejectionNote string) error
	DeleteAchievement(id string) error
}

// achievementRepositoryImpl adalah implementasi dari AchievementRepository
type achievementRepositoryImpl struct {
	db *sql.DB
}

// NewAchievementRepository membuat instance repository achievement baru
func NewAchievementRepository(db *sql.DB) AchievementRepository {
	return &achievementRepositoryImpl{db: db}
}

// CreateAchievement membuat achievement baru
func (r *achievementRepositoryImpl) CreateAchievement(achievement *model.Achievement) error {
	detailsJSON, _ := json.Marshal(achievement.Details)

	query := `
		INSERT INTO achievements (id, student_id, achievement_type, title, description, details, tags, points, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`

	_, err := r.db.Exec(query, achievement.ID, achievement.StudentID, achievement.AchievementType,
		achievement.Title, achievement.Description, detailsJSON, pq.Array(achievement.Tags),
		achievement.Points, achievement.Status)
	return err
}

// GetAchievementByID mengambil achievement berdasarkan ID
func (r *achievementRepositoryImpl) GetAchievementByID(id string) (*model.Achievement, error) {
	query := `
		SELECT id, student_id, achievement_type, title, description, details, tags, points, 
		       status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievements 
		WHERE id = $1
	`

	achievement := &model.Achievement{}
	var detailsJSON []byte
	var tags pq.StringArray

	err := r.db.QueryRow(query, id).Scan(
		&achievement.ID, &achievement.StudentID, &achievement.AchievementType,
		&achievement.Title, &achievement.Description, &detailsJSON, &tags, &achievement.Points,
		&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
		&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal(detailsJSON, &achievement.Details)
	achievement.Tags = []string(tags)

	return achievement, nil
}

// GetAchievementsByStudentID mengambil semua achievements milik student
func (r *achievementRepositoryImpl) GetAchievementsByStudentID(studentID string) ([]*model.Achievement, error) {
	query := `
		SELECT id, student_id, achievement_type, title, description, details, tags, points,
		       status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievements
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*model.Achievement
	for rows.Next() {
		achievement := &model.Achievement{}
		var detailsJSON []byte
		var tags pq.StringArray

		err := rows.Scan(
			&achievement.ID, &achievement.StudentID, &achievement.AchievementType,
			&achievement.Title, &achievement.Description, &detailsJSON, &tags, &achievement.Points,
			&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
			&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(detailsJSON, &achievement.Details)
		achievement.Tags = []string(tags)
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

// GetAchievementsByStatus mengambil achievements berdasarkan status
func (r *achievementRepositoryImpl) GetAchievementsByStatus(status string) ([]*model.Achievement, error) {
	query := `
		SELECT id, student_id, achievement_type, title, description, details, tags, points,
		       status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievements
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*model.Achievement
	for rows.Next() {
		achievement := &model.Achievement{}
		var detailsJSON []byte
		var tags pq.StringArray

		err := rows.Scan(
			&achievement.ID, &achievement.StudentID, &achievement.AchievementType,
			&achievement.Title, &achievement.Description, &detailsJSON, &tags, &achievement.Points,
			&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
			&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(detailsJSON, &achievement.Details)
		achievement.Tags = []string(tags)
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

// GetAllAchievements mengambil semua achievements dengan pagination
func (r *achievementRepositoryImpl) GetAllAchievements(page, pageSize int) ([]*model.Achievement, int, error) {
	offset := (page - 1) * pageSize

	// Hitung total items
	countQuery := `SELECT COUNT(*) FROM achievements`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data dengan pagination
	query := `
		SELECT id, student_id, achievement_type, title, description, details, tags, points,
		       status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievements
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var achievements []*model.Achievement
	for rows.Next() {
		achievement := &model.Achievement{}
		var detailsJSON []byte
		var tags pq.StringArray

		err := rows.Scan(
			&achievement.ID, &achievement.StudentID, &achievement.AchievementType,
			&achievement.Title, &achievement.Description, &detailsJSON, &tags, &achievement.Points,
			&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
			&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		json.Unmarshal(detailsJSON, &achievement.Details)
		achievement.Tags = []string(tags)
		achievements = append(achievements, achievement)
	}

	return achievements, totalItems, nil
}

// UpdateAchievement mengubah data achievement
func (r *achievementRepositoryImpl) UpdateAchievement(achievement *model.Achievement) error {
	detailsJSON, _ := json.Marshal(achievement.Details)

	query := `
		UPDATE achievements
		SET achievement_type = $1, title = $2, description = $3, details = $4, tags = $5, 
		    points = $6, updated_at = NOW()
		WHERE id = $7
	`

	_, err := r.db.Exec(query, achievement.AchievementType, achievement.Title, achievement.Description,
		detailsJSON, pq.Array(achievement.Tags), achievement.Points, achievement.ID)
	return err
}

// SubmitAchievementForVerification mengubah status dari draft ke submitted
func (r *achievementRepositoryImpl) SubmitAchievementForVerification(id string) error {
	query := `
		UPDATE achievements
		SET status = $1, submitted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = $3
	`

	_, err := r.db.Exec(query, "submitted", id, "draft")
	return err
}

// VerifyAchievement mengubah status menjadi verified dan mencatat dosen yang memverifikasi
func (r *achievementRepositoryImpl) VerifyAchievement(id string, verifiedBy string) error {
	query := `
		UPDATE achievements
		SET status = $1, verified_at = NOW(), verified_by = $2, updated_at = NOW()
		WHERE id = $3 AND status = $4
	`

	_, err := r.db.Exec(query, "verified", verifiedBy, id, "submitted")
	return err
}

// RejectAchievement mengubah status menjadi rejected dengan catatan penolakan
func (r *achievementRepositoryImpl) RejectAchievement(id string, verifiedBy string, rejectionNote string) error {
	query := `
		UPDATE achievements
		SET status = $1, verified_at = NOW(), verified_by = $2, rejection_note = $3, updated_at = NOW()
		WHERE id = $4 AND status = $5
	`

	_, err := r.db.Exec(query, "rejected", verifiedBy, rejectionNote, id, "submitted")
	return err
}

// DeleteAchievement menghapus achievement (hanya bisa delete yang draft)
func (r *achievementRepositoryImpl) DeleteAchievement(id string) error {
	query := `
		DELETE FROM achievements WHERE id = $1 AND status = $2
	`

	_, err := r.db.Exec(query, id, "draft")
	return err
}
