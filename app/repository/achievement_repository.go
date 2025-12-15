package repository

import (
	"context"
	"database/sql"
	"time"
	"uas_be/app/model"
	"uas_be/database"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AchievementRepository adalah interface untuk akses data achievement dari database
type AchievementRepository interface {
	CreateAchievement(achievement *model.Achievement, studentID string) (*model.AchievementWithReference, error)
	GetAchievementByID(id string) (*model.AchievementWithReference, error)
	GetAchievementsByStudentID(studentID string) ([]*model.AchievementWithReference, error)
	GetAchievementsByStatus(status string) ([]*model.AchievementWithReference, error)
	GetAllAchievements(page, pageSize int) ([]*model.AchievementWithReference, int, error)
	UpdateAchievement(referenceID string, achievement *model.Achievement) error
	SubmitAchievementForVerification(id string) error
	VerifyAchievement(id string, verifiedBy string) error
	RejectAchievement(id string, verifiedBy string, rejectionNote string) error
	DeleteAchievement(id string) error

	CreateAchievementHistory(history *model.AchievementHistory) error
	GetAchievementHistory(achievementID string) ([]*model.AchievementHistory, error)
	CreateAttachment(attachment *model.AchievementAttachment) error
	GetAttachmentsByAchievementID(achievementID string) ([]*model.AchievementAttachment, error)
}

// achievementRepositoryImpl adalah implementasi dari AchievementRepository
type achievementRepositoryImpl struct {
	db              *sql.DB
	mongoCollection *mongo.Collection // Add MongoDB collection
}

// NewAchievementRepository membuat instance repository achievement baru
func NewAchievementRepository(db *sql.DB) AchievementRepository {
	mongoDB := database.GetMongoDB()
	var collection *mongo.Collection
	if mongoDB != nil {
		collection = mongoDB.Collection("achievements")
	}

	return &achievementRepositoryImpl{
		db:              db,
		mongoCollection: collection,
	}
}

// CreateAchievement menyimpan achievement ke MongoDB dan reference ke PostgreSQL
func (r *achievementRepositoryImpl) CreateAchievement(achievement *model.Achievement, studentID string) (*model.AchievementWithReference, error) {
	ctx := context.Background()

	// 1. Insert achievement data to MongoDB
	achievement.StudentID = studentID
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := r.mongoCollection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	mongoID := result.InsertedID.(primitive.ObjectID).Hex()

	// 2. Create reference in PostgreSQL
	referenceID := uuid.New().String()
	query := `
		INSERT INTO achievement_references (id, student_id, mongo_achievement_id, achievement_title, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err = r.db.Exec(query, referenceID, studentID, mongoID, achievement.Title, model.AchievementStatusDraft)
	if err != nil {
		// Rollback: delete from MongoDB if PostgreSQL insert fails
		r.mongoCollection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		return nil, err
	}

	// 3. Return combined result
	return &model.AchievementWithReference{
		Achievement: *achievement,
		Status:      model.AchievementStatusDraft,
	}, nil
}

// GetAchievementByID mengambil achievement dari MongoDB dan reference dari PostgreSQL
func (r *achievementRepositoryImpl) GetAchievementByID(referenceID string) (*model.AchievementWithReference, error) {
	ctx := context.Background()

	// 1. Get reference from PostgreSQL
	query := `
		SELECT id, student_id, mongo_achievement_id, achievement_title, status, 
		       submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references 
		WHERE id = $1
	`

	var ref model.AchievementReference
	err := r.db.QueryRow(query, referenceID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.AchievementTitle,
		&ref.Status, &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
		&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 2. Get achievement data from MongoDB
	mongoObjID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	err = r.mongoCollection.FindOne(ctx, bson.M{"_id": mongoObjID}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	// 3. Combine both
	return &model.AchievementWithReference{
		Achievement:   achievement,
		Status:        ref.Status,
		SubmittedAt:   ref.SubmittedAt,
		VerifiedAt:    ref.VerifiedAt,
		VerifiedBy:    ref.VerifiedBy,
		RejectionNote: ref.RejectionNote,
	}, nil
}

func (r *achievementRepositoryImpl) GetAchievementsByStudentID(studentID string) ([]*model.AchievementWithReference, error) {
	ctx := context.Background()

	// 1. Get all references from PostgreSQL
	query := `
		SELECT id, student_id, mongo_achievement_id, achievement_title, status,
		       submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.AchievementWithReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.AchievementTitle,
			&ref.Status, &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// 2. Get achievement data from MongoDB
		mongoObjID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		var achievement model.Achievement
		err = r.mongoCollection.FindOne(ctx, bson.M{"_id": mongoObjID}).Decode(&achievement)
		if err != nil {
			continue
		}

		// 3. Combine
		results = append(results, &model.AchievementWithReference{
			Achievement:   achievement,
			Status:        ref.Status,
			SubmittedAt:   ref.SubmittedAt,
			VerifiedAt:    ref.VerifiedAt,
			VerifiedBy:    ref.VerifiedBy,
			RejectionNote: ref.RejectionNote,
		})
	}

	return results, nil
}

func (r *achievementRepositoryImpl) GetAchievementsByStatus(status string) ([]*model.AchievementWithReference, error) {
	ctx := context.Background()

	query := `
		SELECT id, student_id, mongo_achievement_id, achievement_title, status,
		       submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.AchievementWithReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.AchievementTitle,
			&ref.Status, &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			continue
		}

		mongoObjID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		var achievement model.Achievement
		err = r.mongoCollection.FindOne(ctx, bson.M{"_id": mongoObjID}).Decode(&achievement)
		if err != nil {
			continue
		}

		results = append(results, &model.AchievementWithReference{
			Achievement:   achievement,
			Status:        ref.Status,
			SubmittedAt:   ref.SubmittedAt,
			VerifiedAt:    ref.VerifiedAt,
			VerifiedBy:    ref.VerifiedBy,
			RejectionNote: ref.RejectionNote,
		})
	}

	return results, nil
}

func (r *achievementRepositoryImpl) GetAllAchievements(page, pageSize int) ([]*model.AchievementWithReference, int, error) {
	ctx := context.Background()
	offset := (page - 1) * pageSize

	// Count total
	countQuery := `SELECT COUNT(*) FROM achievement_references`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated references
	query := `
		SELECT id, student_id, mongo_achievement_id, achievement_title, status,
		       submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*model.AchievementWithReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.AchievementTitle,
			&ref.Status, &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			continue
		}

		mongoObjID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		var achievement model.Achievement
		err = r.mongoCollection.FindOne(ctx, bson.M{"_id": mongoObjID}).Decode(&achievement)
		if err != nil {
			continue
		}

		results = append(results, &model.AchievementWithReference{
			Achievement:   achievement,
			Status:        ref.Status,
			SubmittedAt:   ref.SubmittedAt,
			VerifiedAt:    ref.VerifiedAt,
			VerifiedBy:    ref.VerifiedBy,
			RejectionNote: ref.RejectionNote,
		})
	}

	return results, totalItems, nil
}

func (r *achievementRepositoryImpl) UpdateAchievement(referenceID string, achievement *model.Achievement) error {
	ctx := context.Background()

	// Get reference to find MongoDB ID
	var mongoID string
	err := r.db.QueryRow("SELECT mongo_achievement_id FROM achievement_references WHERE id = $1", referenceID).Scan(&mongoID)
	if err != nil {
		return err
	}

	mongoObjID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	// Update MongoDB
	achievement.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"achievement_type": achievement.AchievementType,
			"title":            achievement.Title,
			"description":      achievement.Description,
			"details":          achievement.Details,
			"tags":             achievement.Tags,
			"points":           achievement.Points,
			"updated_at":       achievement.UpdatedAt,
		},
	}

	_, err = r.mongoCollection.UpdateOne(ctx, bson.M{"_id": mongoObjID}, update)
	if err != nil {
		return err
	}

	// Update title in PostgreSQL reference
	_, err = r.db.Exec("UPDATE achievement_references SET achievement_title = $1, updated_at = NOW() WHERE id = $2", achievement.Title, referenceID)
	return err
}

func (r *achievementRepositoryImpl) SubmitAchievementForVerification(id string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, submitted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = $3
	`
	_, err := r.db.Exec(query, "submitted", id, "draft")
	return err
}

func (r *achievementRepositoryImpl) VerifyAchievement(id string, verifiedBy string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_at = NOW(), verified_by = $2, updated_at = NOW()
		WHERE id = $3 AND status = $4
	`
	_, err := r.db.Exec(query, "verified", verifiedBy, id, "submitted")
	return err
}

func (r *achievementRepositoryImpl) RejectAchievement(id string, verifiedBy string, rejectionNote string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_at = NOW(), verified_by = $2, rejection_note = $3, updated_at = NOW()
		WHERE id = $4 AND status = $5
	`
	_, err := r.db.Exec(query, "rejected", verifiedBy, rejectionNote, id, "submitted")
	return err
}

func (r *achievementRepositoryImpl) DeleteAchievement(id string) error {
	ctx := context.Background()

	// Get MongoDB ID
	var mongoID string
	var status string
	err := r.db.QueryRow("SELECT mongo_achievement_id, status FROM achievement_references WHERE id = $1", id).Scan(&mongoID, &status)
	if err != nil {
		return err
	}

	if status != "draft" {
		return nil // Only draft can be deleted
	}

	// Delete from PostgreSQL
	_, err = r.db.Exec("DELETE FROM achievement_references WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Delete from MongoDB
	mongoObjID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = r.mongoCollection.DeleteOne(ctx, bson.M{"_id": mongoObjID})
	return err
}

// CreateAchievementHistory menyimpan history perubahan status achievement
func (r *achievementRepositoryImpl) CreateAchievementHistory(history *model.AchievementHistory) error {
	query := `
		INSERT INTO achievement_history (id, achievement_id, old_status, new_status, changed_by, note, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Exec(query, history.ID, history.AchievementID, history.OldStatus, history.NewStatus, history.ChangedBy, history.Note)
	return err
}

// GetAchievementHistory mengambil riwayat perubahan achievement
func (r *achievementRepositoryImpl) GetAchievementHistory(achievementID string) ([]*model.AchievementHistory, error) {
	query := `
		SELECT ah.id, ah.achievement_id, ah.old_status, ah.new_status, ah.changed_by, u.full_name, ah.note, ah.created_at
		FROM achievement_history ah
		JOIN users u ON ah.changed_by = u.id
		WHERE ah.achievement_id = $1
		ORDER BY ah.created_at DESC
	`

	rows, err := r.db.Query(query, achievementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*model.AchievementHistory
	for rows.Next() {
		history := &model.AchievementHistory{}
		err := rows.Scan(&history.ID, &history.AchievementID, &history.OldStatus, &history.NewStatus, &history.ChangedBy, &history.ChangedByName, &history.Note, &history.CreatedAt)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}

// CreateAttachment menyimpan attachment file untuk achievement
func (r *achievementRepositoryImpl) CreateAttachment(attachment *model.AchievementAttachment) error {
	query := `
		INSERT INTO achievement_attachments (id, achievement_id, file_name, file_url, file_size, file_type, uploaded_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`
	_, err := r.db.Exec(query, attachment.ID, attachment.AchievementID, attachment.FileName, attachment.FileURL, attachment.FileSize, attachment.FileType, attachment.UploadedBy)
	return err
}

// GetAttachmentsByAchievementID mengambil semua attachment dari achievement
func (r *achievementRepositoryImpl) GetAttachmentsByAchievementID(achievementID string) ([]*model.AchievementAttachment, error) {
	query := `
		SELECT id, achievement_id, file_name, file_url, file_size, file_type, uploaded_by, created_at
		FROM achievement_attachments
		WHERE achievement_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, achievementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*model.AchievementAttachment
	for rows.Next() {
		attachment := &model.AchievementAttachment{}
		err := rows.Scan(&attachment.ID, &attachment.AchievementID, &attachment.FileName, &attachment.FileURL, &attachment.FileSize, &attachment.FileType, &attachment.UploadedBy, &attachment.CreatedAt)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}
