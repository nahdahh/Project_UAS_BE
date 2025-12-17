package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	GetAchievementsWithFilters(page, pageSize int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.AchievementWithReference, int, error)
	UpdateAchievement(referenceID string, achievement *model.Achievement) error
	SubmitAchievementForVerification(id string) error
	VerifyAchievement(id string, verifiedBy string) error
	RejectAchievement(id string, verifiedBy string, rejectionNote string) error
	DeleteAchievement(id string) error

	CreateAchievementHistory(history *model.AchievementHistory) error
	GetAchievementHistory(achievementID string) ([]*model.AchievementHistory, error)
	CreateAttachment(attachment *model.AchievementAttachment) error
	GetAttachmentsByAchievementID(achievementID string) ([]*model.AchievementAttachment, error)

	GetAchievementStatsByPeriod(startDate, endDate time.Time, role, userID string) (map[string]interface{}, error)
	GetAchievementStatsByType(role, userID string) (map[string]interface{}, error)
	GetTopStudents(limit int) ([]*model.StudentStats, error)
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
		return nil, fmt.Errorf("failed to insert achievement to MongoDB: %w", err)
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
		return nil, fmt.Errorf("failed to create achievement reference in PostgreSQL: %w", err)
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
		       submitted_at, verified_at, verified_by, rejection_note, deleted_at, created_at, updated_at
		FROM achievement_references 
		WHERE id = $1
	`

	var ref model.AchievementReference
	err := r.db.QueryRow(query, referenceID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.AchievementTitle,
		&ref.Status, &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
		&ref.RejectionNote, &ref.DeletedAt, &ref.CreatedAt, &ref.UpdatedAt,
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
		       submitted_at, verified_at, verified_by, rejection_note, deleted_at, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1 AND status != $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, studentID, model.AchievementStatusDeleted)
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
			&ref.RejectionNote, &ref.DeletedAt, &ref.CreatedAt, &ref.UpdatedAt,
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
		       submitted_at, verified_at, verified_by, rejection_note, deleted_at, created_at, updated_at
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
			&ref.RejectionNote, &ref.DeletedAt, &ref.CreatedAt, &ref.UpdatedAt,
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
	countQuery := `SELECT COUNT(*) FROM achievement_references WHERE status != $1`
	var totalItems int
	err := r.db.QueryRow(countQuery, model.AchievementStatusDeleted).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated references
	query := `
		SELECT id, student_id, mongo_achievement_id, achievement_title, status,
		       submitted_at, verified_at, verified_by, rejection_note, deleted_at, created_at, updated_at
		FROM achievement_references
		WHERE status != $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, model.AchievementStatusDeleted, pageSize, offset)
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
			&ref.RejectionNote, &ref.DeletedAt, &ref.CreatedAt, &ref.UpdatedAt,
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

func (r *achievementRepositoryImpl) GetAchievementsWithFilters(page, pageSize int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.AchievementWithReference, int, error) {
	ctx := context.Background()
	offset := (page - 1) * pageSize

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	whereClauses = append(whereClauses, fmt.Sprintf("status != $%d", argCounter))
	args = append(args, model.AchievementStatusDeleted)
	argCounter++

	if status, ok := filters["status"].(string); ok && status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, status)
		argCounter++
	}

	if studentID, ok := filters["student_id"].(string); ok && studentID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("student_id = $%d", argCounter))
		args = append(args, studentID)
		argCounter++
	}

	if studentIDs, ok := filters["student_ids"].([]string); ok && len(studentIDs) > 0 {
		placeholders := make([]string, len(studentIDs))
		for i := range studentIDs {
			placeholders[i] = fmt.Sprintf("$%d", argCounter)
			args = append(args, studentIDs[i])
			argCounter++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("student_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argCounter))
		args = append(args, startDate)
		argCounter++
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argCounter))
		args = append(args, endDate)
		argCounter++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM achievement_references WHERE %s", whereClause)
	var totalItems int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" || (sortOrder != "ASC" && sortOrder != "DESC") {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT id, student_id, mongo_achievement_id, achievement_title, status,
		       submitted_at, verified_at, verified_by, rejection_note, deleted_at, created_at, updated_at
		FROM achievement_references
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortOrder, argCounter, argCounter+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
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
			&ref.RejectionNote, &ref.DeletedAt, &ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			continue
		}

		mongoObjID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		var achievement model.Achievement
		mongoFilter := bson.M{"_id": mongoObjID}

		if achievementType, ok := filters["achievement_type"].(string); ok && achievementType != "" {
			mongoFilter["achievement_type"] = achievementType
		}

		err = r.mongoCollection.FindOne(ctx, mongoFilter).Decode(&achievement)
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

	var mongoID string
	err := r.db.QueryRow("SELECT mongo_achievement_id FROM achievement_references WHERE id = $1", referenceID).Scan(&mongoID)
	if err != nil {
		return err
	}

	mongoObjID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

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
	query := `
		UPDATE achievement_references
		SET status = $1, deleted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = $3
	`
	_, err := r.db.Exec(query, model.AchievementStatusDeleted, id, model.AchievementStatusDraft)
	return err
}

// CreateAchievementHistory menyimpan history perubahan status achievement
func (r *achievementRepositoryImpl) CreateAchievementHistory(history *model.AchievementHistory) error {
	query := `
		INSERT INTO achievement_history (id, achievement_id, previous_status, new_status, changed_by, notes, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Exec(query, history.ID, history.AchievementID, history.OldStatus, history.NewStatus, history.ChangedBy, history.Note)
	return err
}

// GetAchievementHistory mengambil riwayat perubahan achievement
func (r *achievementRepositoryImpl) GetAchievementHistory(achievementID string) ([]*model.AchievementHistory, error) {
	query := `
		SELECT ah.id, ah.achievement_id, ah.previous_status, ah.new_status, ah.changed_by, u.full_name, ah.notes, ah.changed_at
		FROM achievement_history ah
		JOIN users u ON ah.changed_by = u.id
		WHERE ah.achievement_id = $1
		ORDER BY ah.changed_at DESC
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
		INSERT INTO achievement_attachments (id, achievement_id, file_name, file_path, file_size, file_type, uploaded_by, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`
	_, err := r.db.Exec(query, attachment.ID, attachment.AchievementID, attachment.FileName, attachment.FilePath, attachment.FileSize, attachment.FileType, attachment.UploadedBy)
	return err
}

// GetAttachmentsByAchievementID mengambil semua attachment dari achievement
func (r *achievementRepositoryImpl) GetAttachmentsByAchievementID(achievementID string) ([]*model.AchievementAttachment, error) {
	query := `
		SELECT id, achievement_id, file_name, file_path, file_size, file_type, uploaded_by, uploaded_at
		FROM achievement_attachments
		WHERE achievement_id = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.Query(query, achievementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*model.AchievementAttachment
	for rows.Next() {
		attachment := &model.AchievementAttachment{}
		err := rows.Scan(&attachment.ID, &attachment.AchievementID, &attachment.FileName, &attachment.FilePath, &attachment.FileSize, &attachment.FileType, &attachment.UploadedBy, &attachment.UploadedAt)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// GetAchievementStatsByPeriod mengambil statistik achievement berdasarkan periode waktu
func (r *achievementRepositoryImpl) GetAchievementStatsByPeriod(startDate, endDate time.Time, role, userID string) (map[string]interface{}, error) {
	ctx := context.Background()

	whereClause := "status != $1 AND created_at >= $2 AND created_at <= $3"
	args := []interface{}{model.AchievementStatusDeleted, startDate, endDate}

	if role == "Mahasiswa" {
		whereClause += " AND student_id = $4"
		args = append(args, userID)
	} else if role == "Dosen Wali" {
		whereClause += " AND student_id IN (SELECT id FROM students WHERE advisor_id = $4)"
		args = append(args, userID)
	}

	query := fmt.Sprintf(`
		SELECT id, student_id, mongo_achievement_id, status
		FROM achievement_references
		WHERE %s
	`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusCount := make(map[string]int)
	typeCount := make(map[string]int)
	totalPoints := 0
	var mongoIDs []primitive.ObjectID

	for rows.Next() {
		var id, studentID, mongoID, status string
		if err := rows.Scan(&id, &studentID, &mongoID, &status); err != nil {
			continue
		}

		statusCount[status]++

		objID, err := primitive.ObjectIDFromHex(mongoID)
		if err != nil {
			continue
		}
		mongoIDs = append(mongoIDs, objID)
	}

	for _, objID := range mongoIDs {
		var achievement model.Achievement
		err := r.mongoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&achievement)
		if err != nil {
			continue
		}

		typeCount[achievement.AchievementType]++
		if statusCount["verified"] > 0 {
			totalPoints += achievement.Points
		}
	}

	stats := map[string]interface{}{
		"period": map[string]interface{}{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
		"total_achievements": len(mongoIDs),
		"by_status":          statusCount,
		"by_type":            typeCount,
		"total_points":       totalPoints,
	}

	return stats, nil
}

// GetAchievementStatsByType mengambil statistik achievement berdasarkan jenis achievement
func (r *achievementRepositoryImpl) GetAchievementStatsByType(role, userID string) (map[string]interface{}, error) {
	ctx := context.Background()

	whereClause := "status != $1"
	args := []interface{}{model.AchievementStatusDeleted}

	if role == "Mahasiswa" {
		whereClause += " AND student_id = $2"
		args = append(args, userID)
	} else if role == "Dosen Wali" {
		whereClause += " AND student_id IN (SELECT id FROM students WHERE advisor_id = $2)"
		args = append(args, userID)
	}

	query := fmt.Sprintf(`
		SELECT mongo_achievement_id, status
		FROM achievement_references
		WHERE %s
	`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	typeStats := make(map[string]map[string]int)

	for rows.Next() {
		var mongoID, status string
		if err := rows.Scan(&mongoID, &status); err != nil {
			continue
		}

		objID, err := primitive.ObjectIDFromHex(mongoID)
		if err != nil {
			continue
		}

		var achievement model.Achievement
		err = r.mongoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&achievement)
		if err != nil {
			continue
		}

		if _, ok := typeStats[achievement.AchievementType]; !ok {
			typeStats[achievement.AchievementType] = make(map[string]int)
		}
		typeStats[achievement.AchievementType][status]++
	}

	return map[string]interface{}{
		"stats_by_type": typeStats,
	}, nil
}

// GetTopStudents mengambil top students berdasarkan jumlah achievement yang diverifikasi
func (r *achievementRepositoryImpl) GetTopStudents(limit int) ([]*model.StudentStats, error) {
	ctx := context.Background()

	query := `
		SELECT ar.student_id, s.nim, s.name, COUNT(ar.id) as total_achievements
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE ar.status = $1
		GROUP BY ar.student_id, s.nim, s.name
		ORDER BY total_achievements DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, "verified", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topStudents []*model.StudentStats

	for rows.Next() {
		var studentID, nim, name string
		var totalAchievements int

		if err := rows.Scan(&studentID, &nim, &name, &totalAchievements); err != nil {
			continue
		}

		mongoQuery := `
			SELECT mongo_achievement_id
			FROM achievement_references
			WHERE student_id = $1 AND status = $2
		`

		mongoRows, err := r.db.Query(mongoQuery, studentID, "verified")
		if err != nil {
			continue
		}

		totalPoints := 0
		for mongoRows.Next() {
			var mongoID string
			if err := mongoRows.Scan(&mongoID); err != nil {
				continue
			}

			objID, err := primitive.ObjectIDFromHex(mongoID)
			if err != nil {
				continue
			}

			var achievement model.Achievement
			err = r.mongoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&achievement)
			if err != nil {
				continue
			}

			totalPoints += achievement.Points
		}
		mongoRows.Close()

		topStudents = append(topStudents, &model.StudentStats{
			StudentID:         studentID,
			NIM:               nim,
			Name:              name,
			TotalAchievements: totalAchievements,
			TotalPoints:       totalPoints,
		})
	}

	return topStudents, nil
}
