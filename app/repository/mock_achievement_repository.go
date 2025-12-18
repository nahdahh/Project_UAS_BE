package repository

import (
	"errors"
	"time"
	"uas_be/app/model"

	"github.com/google/uuid"
)

type MockAchievementRepository struct {
	achievements map[string]*model.AchievementWithReference
	histories    map[string][]*model.AchievementHistory
	attachments  map[string][]*model.AchievementAttachment
}

func NewMockAchievementRepository() *MockAchievementRepository {
	return &MockAchievementRepository{
		achievements: make(map[string]*model.AchievementWithReference),
		histories:    make(map[string][]*model.AchievementHistory),
		attachments:  make(map[string][]*model.AchievementAttachment),
	}
}

func (m *MockAchievementRepository) Create(achievement *model.Achievement, studentID string) (*model.AchievementWithReference, error) {
	achievement.StudentID = studentID
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	achWithRef := &model.AchievementWithReference{
		Achievement: *achievement,
		Status:      model.AchievementStatusDraft,
	}
	// Use studentID as key for simplicity in tests
	m.achievements[studentID] = achWithRef
	return achWithRef, nil
}

func (m *MockAchievementRepository) CreateAchievement(achievement *model.Achievement, studentID string) (*model.AchievementWithReference, error) {
	return m.Create(achievement, studentID)
}

func (m *MockAchievementRepository) GetAchievementByID(id string) (*model.AchievementWithReference, error) {
	if achievement, exists := m.achievements[id]; exists {
		return achievement, nil
	}
	return nil, nil
}

func (m *MockAchievementRepository) GetAchievementsByStudentID(studentID string) ([]*model.AchievementWithReference, error) {
	var achievements []*model.AchievementWithReference
	for _, achievement := range m.achievements {
		if achievement.StudentID == studentID {
			achievements = append(achievements, achievement)
		}
	}
	return achievements, nil
}

func (m *MockAchievementRepository) GetAchievementsByStatus(status string) ([]*model.AchievementWithReference, error) {
	var achievements []*model.AchievementWithReference
	for _, achievement := range m.achievements {
		if achievement.Status == status {
			achievements = append(achievements, achievement)
		}
	}
	return achievements, nil
}

func (m *MockAchievementRepository) GetAllAchievements(page, pageSize int) ([]*model.AchievementWithReference, int, error) {
	var achievements []*model.AchievementWithReference
	for _, achievement := range m.achievements {
		if achievement.Status != model.AchievementStatusDeleted {
			achievements = append(achievements, achievement)
		}
	}
	return achievements, len(achievements), nil
}

func (m *MockAchievementRepository) GetAchievementsWithFilters(page, pageSize int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.AchievementWithReference, int, error) {
	var achievements []*model.AchievementWithReference
	for _, achievement := range m.achievements {
		if achievement.Status != model.AchievementStatusDeleted {
			// Apply filters
			if status, ok := filters["status"].(string); ok && status != "" {
				if achievement.Status != status {
					continue
				}
			}
			if studentID, ok := filters["student_id"].(string); ok && studentID != "" {
				if achievement.StudentID != studentID {
					continue
				}
			}
			if studentIDs, ok := filters["student_ids"].([]string); ok {
				found := false
				for _, sid := range studentIDs {
					if achievement.StudentID == sid {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			achievements = append(achievements, achievement)
		}
	}
	return achievements, len(achievements), nil
}

func (m *MockAchievementRepository) UpdateAchievement(id string, achievement *model.Achievement) error {
	if ach, exists := m.achievements[id]; exists {
		ach.Achievement = *achievement
		ach.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("achievement tidak ditemukan")
}

func (m *MockAchievementRepository) Submit(id string) error {
	if ach, exists := m.achievements[id]; exists {
		ach.Status = model.AchievementStatusSubmitted
		now := time.Now()
		ach.SubmittedAt = &now
		return nil
	}
	return errors.New("achievement tidak ditemukan")
}

func (m *MockAchievementRepository) SubmitAchievementForVerification(id string) error {
	return m.Submit(id)
}

func (m *MockAchievementRepository) VerifyAchievement(id string, verifiedBy string) error {
	if ach, exists := m.achievements[id]; exists {
		if ach.Status != model.AchievementStatusSubmitted {
			return errors.New("hanya prestasi submitted yang bisa diverify")
		}
		ach.Status = model.AchievementStatusVerified
		now := time.Now()
		ach.VerifiedAt = &now
		ach.VerifiedBy = &verifiedBy
		return nil
	}
	return errors.New("achievement tidak ditemukan")
}

func (m *MockAchievementRepository) RejectAchievement(id string, verifiedBy string, rejectionNote string) error {
	if ach, exists := m.achievements[id]; exists {
		if ach.Status != model.AchievementStatusSubmitted {
			return errors.New("hanya prestasi submitted yang bisa direject")
		}
		ach.Status = model.AchievementStatusRejected
		now := time.Now()
		ach.VerifiedAt = &now
		ach.VerifiedBy = &verifiedBy
		ach.RejectionNote = &rejectionNote
		return nil
	}
	return errors.New("achievement tidak ditemukan")
}

func (m *MockAchievementRepository) DeleteAchievement(id string) error {
	if ach, exists := m.achievements[id]; exists {
		if ach.Status != model.AchievementStatusDraft {
			return errors.New("hanya prestasi draft yang bisa dihapus")
		}
		ach.Status = model.AchievementStatusDeleted
		now := time.Now()
		ach.Achievement.UpdatedAt = now
		return nil
	}
	return errors.New("achievement tidak ditemukan")
}

func (m *MockAchievementRepository) CreateAchievementHistory(history *model.AchievementHistory) error {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}
	history.CreatedAt = time.Now()
	m.histories[history.AchievementID] = append(m.histories[history.AchievementID], history)
	return nil
}

func (m *MockAchievementRepository) GetAchievementHistory(achievementID string) ([]*model.AchievementHistory, error) {
	if histories, exists := m.histories[achievementID]; exists {
		return histories, nil
	}
	return []*model.AchievementHistory{}, nil
}

func (m *MockAchievementRepository) CreateAttachment(attachment *model.AchievementAttachment) error {
	if attachment.ID == "" {
		attachment.ID = uuid.New().String()
	}
	attachment.UploadedAt = time.Now()
	m.attachments[attachment.AchievementID] = append(m.attachments[attachment.AchievementID], attachment)
	return nil
}

func (m *MockAchievementRepository) GetAttachmentsByAchievementID(achievementID string) ([]*model.AchievementAttachment, error) {
	if attachments, exists := m.attachments[achievementID]; exists {
		return attachments, nil
	}
	return []*model.AchievementAttachment{}, nil
}

func (m *MockAchievementRepository) GetAchievementStatsByPeriod(startDate, endDate time.Time, role, userID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"total": len(m.achievements),
	}, nil
}

func (m *MockAchievementRepository) GetAchievementStatsByType(role, userID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"academic":      0,
		"competition":   0,
		"organization":  0,
		"publication":   0,
		"certification": 0,
		"other":         0,
	}, nil
}

func (m *MockAchievementRepository) GetTopStudents(limit int) ([]*model.StudentStats, error) {
	return []*model.StudentStats{}, nil
}
