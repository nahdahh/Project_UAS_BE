package repository

import (
	"database/sql"
	"uas_be/app/model"
)

// AchievementRepository adalah interface untuk akses data achievement dari database
type AchievementRepository interface {
	// CreateAchievement membuat achievement baru
	CreateAchievement(achievement *model.AchievementReference) error

	// GetAchievementByID mengambil achievement berdasarkan ID
	GetAchievementByID(id string) (*model.AchievementReference, error)

	// GetAchievementsByStudentID mengambil semua achievement milik student
	GetAchievementsByStudentID(studentID string) ([]*model.AchievementReference, error)

	// GetAchievementsByStatus mengambil achievement berdasarkan status (draft, submitted, verified, rejected)
	GetAchievementsByStatus(status string) ([]*model.AchievementReference, error)

	// GetAllAchievements mengambil semua achievement dengan pagination
	GetAllAchievements(page, pageSize int) ([]*model.AchievementReference, int, error)

	// UpdateAchievement mengubah data achievement
	UpdateAchievement(achievement *model.AchievementReference) error

	// SubmitAchievementForVerification mengubah status achievement dari draft menjadi submitted
	SubmitAchievementForVerification(id string) error

	// VerifyAchievement mengubah status achievement menjadi verified dan mencatat dosen yang memverifikasi
	VerifyAchievement(id string, verifiedBy string) error

	// RejectAchievement mengubah status achievement menjadi rejected dengan catatan penolakan
	RejectAchievement(id string, verifiedBy string, rejectionNote string) error

	// DeleteAchievement menghapus achievement (soft delete)
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

// CreateAchievement membuat achievement baru di database PostgreSQL
func (r *achievementRepositoryImpl) CreateAchievement(achievement *model.AchievementReference) error {
	query := `
		INSERT INTO achievement_references (id, student_id, achievement_title, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := r.db.Exec(query, achievement.ID, achievement.StudentID, achievement.AchievementTitle, achievement.Status)
	return err
}

// GetAchievementByID mengambil achievement berdasarkan ID
func (r *achievementRepositoryImpl) GetAchievementByID(id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, achievement_title, status, submitted_at, verified_at, verified_by, 
		       rejection_note, created_at, updated_at
		FROM achievement_references 
		WHERE id = $1
	`

	achievement := &model.AchievementReference{}
	err := r.db.QueryRow(query, id).Scan(
		&achievement.ID, &achievement.StudentID, &achievement.AchievementTitle,
		&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
		&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return achievement, nil
}

// GetAchievementsByStudentID mengambil semua achievement milik student
func (r *achievementRepositoryImpl) GetAchievementsByStudentID(studentID string) ([]*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, achievement_title, status, submitted_at, verified_at, verified_by,
		       rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*model.AchievementReference
	for rows.Next() {
		achievement := &model.AchievementReference{}
		err := rows.Scan(
			&achievement.ID, &achievement.StudentID, &achievement.AchievementTitle,
			&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
			&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

// GetAchievementsByStatus mengambil achievement berdasarkan status
func (r *achievementRepositoryImpl) GetAchievementsByStatus(status string) ([]*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, achievement_title, status, submitted_at, verified_at, verified_by,
		       rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*model.AchievementReference
	for rows.Next() {
		achievement := &model.AchievementReference{}
		err := rows.Scan(
			&achievement.ID, &achievement.StudentID, &achievement.AchievementTitle,
			&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
			&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

// GetAllAchievements mengambil semua achievement dengan pagination
func (r *achievementRepositoryImpl) GetAllAchievements(page, pageSize int) ([]*model.AchievementReference, int, error) {
	offset := (page - 1) * pageSize

	// Hitung total items
	countQuery := `SELECT COUNT(*) FROM achievement_references`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data dengan pagination
	query := `
		SELECT id, student_id, achievement_title, status, submitted_at, verified_at, verified_by,
		       rejection_note, created_at, updated_at
		FROM achievement_references
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var achievements []*model.AchievementReference
	for rows.Next() {
		achievement := &model.AchievementReference{}
		err := rows.Scan(
			&achievement.ID, &achievement.StudentID, &achievement.AchievementTitle,
			&achievement.Status, &achievement.SubmittedAt, &achievement.VerifiedAt,
			&achievement.VerifiedBy, &achievement.RejectionNote, &achievement.CreatedAt, &achievement.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		achievements = append(achievements, achievement)
	}

	return achievements, totalItems, nil
}

// UpdateAchievement mengubah data achievement
func (r *achievementRepositoryImpl) UpdateAchievement(achievement *model.AchievementReference) error {
	query := `
		UPDATE achievement_references
		SET achievement_title = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(query, achievement.AchievementTitle, achievement.ID)
	return err
}

// SubmitAchievementForVerification mengubah status dari draft ke submitted
func (r *achievementRepositoryImpl) SubmitAchievementForVerification(id string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, submitted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = $3
	`

	_, err := r.db.Exec(query, "submitted", id, "draft")
	return err
}

// VerifyAchievement mengubah status menjadi verified dan mencatat dosen yang memverifikasi
func (r *achievementRepositoryImpl) VerifyAchievement(id string, verifiedBy string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_at = NOW(), verified_by = $2, updated_at = NOW()
		WHERE id = $3 AND status = $4
	`

	_, err := r.db.Exec(query, "verified", verifiedBy, id, "submitted")
	return err
}

// RejectAchievement mengubah status menjadi rejected dengan catatan penolakan
func (r *achievementRepositoryImpl) RejectAchievement(id string, verifiedBy string, rejectionNote string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_at = NOW(), verified_by = $2, rejection_note = $3, updated_at = NOW()
		WHERE id = $4 AND status = $5
	`

	_, err := r.db.Exec(query, "rejected", verifiedBy, rejectionNote, id, "submitted")
	return err
}

// DeleteAchievement menghapus achievement (soft delete dengan mengubah status)
func (r *achievementRepositoryImpl) DeleteAchievement(id string) error {
	query := `
		DELETE FROM achievement_references WHERE id = $1 AND status = $2
	`

	_, err := r.db.Exec(query, id, "draft")
	return err
}
