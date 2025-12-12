package service

import (
	"strconv"
	"time"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementService interface {
	GetAllAchievements(c *fiber.Ctx) error
	GetAchievementDetail(c *fiber.Ctx) error
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
	SubmitAchievement(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error
	GetAdviseeAchievements(c *fiber.Ctx) error
}

type achievementServiceImpl struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
) AchievementService {
	return &achievementServiceImpl{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
	}
}

func (s *achievementServiceImpl) GetAllAchievements(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	achievements, total, err := s.achievementRepo.GetAllAchievements(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil achievements",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievements berhasil diambil",
		Data: map[string]interface{}{
			"achievements": achievements,
			"pagination": map[string]interface{}{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + pageSize - 1) / pageSize,
			},
		},
	})
}

func (s *achievementServiceImpl) GetAchievementDetail(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	studentID := c.Locals("userID").(string)

	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil achievement",
		})
	}
	if achievement == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi tidak ditemukan",
		})
	}

	role := c.Locals("role").(string)
	if role == "Mahasiswa" && achievement.StudentID != studentID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil diambil",
		Data:    achievement,
	})
}

func (s *achievementServiceImpl) CreateAchievement(c *fiber.Ctx) error {
	studentID := c.Locals("userID").(string)

	var req model.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	// Validate student exists
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil || student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	// Create achievement
	achievement := &model.Achievement{
		ID:              uuid.New().String(),
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		Status:          model.AchievementStatusDraft,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.achievementRepo.CreateAchievement(achievement); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal membuat achievement",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil dibuat",
		Data:    achievement,
	})
}

func (s *achievementServiceImpl) UpdateAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	studentID := c.Locals("userID").(string)

	var req model.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil || achievement == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi tidak ditemukan",
		})
	}

	if achievement.StudentID != studentID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi bukan milik anda",
		})
	}

	if achievement.Status != model.AchievementStatusDraft {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "hanya prestasi draft yang bisa diupdate",
		})
	}

	// Update fields if provided
	if req.AchievementType != nil {
		achievement.AchievementType = *req.AchievementType
	}
	if req.Title != nil {
		achievement.Title = *req.Title
	}
	if req.Description != nil {
		achievement.Description = *req.Description
	}
	if req.Details != nil {
		achievement.Details = *req.Details
	}
	if req.Tags != nil {
		achievement.Tags = *req.Tags
	}
	if req.Points != nil {
		achievement.Points = *req.Points
	}

	if err := s.achievementRepo.UpdateAchievement(achievement); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal update achievement",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil diupdate",
		Data:    achievement,
	})
}

func (s *achievementServiceImpl) SubmitAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	studentID := c.Locals("userID").(string)

	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil || achievement == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi tidak ditemukan",
		})
	}

	if achievement.StudentID != studentID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi bukan milik anda",
		})
	}

	if achievement.Status != model.AchievementStatusDraft {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "hanya prestasi draft yang bisa disubmit",
		})
	}

	if err := s.achievementRepo.SubmitAchievementForVerification(achievementID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal submit achievement",
		})
	}

	// Refresh data
	achievement, _ = s.achievementRepo.GetAchievementByID(achievementID)
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil disubmit",
		Data:    achievement,
	})
}

func (s *achievementServiceImpl) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	verifiedBy := c.Locals("userID").(string)

	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil || achievement == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi tidak ditemukan",
		})
	}

	if achievement.Status != model.AchievementStatusSubmitted {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "hanya prestasi submitted yang bisa diverify",
		})
	}

	student, _ := s.studentRepo.GetStudentByID(achievement.StudentID)
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	if student.AdvisorID != verifiedBy {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "anda bukan advisor dari student ini",
		})
	}

	if err := s.achievementRepo.VerifyAchievement(achievementID, verifiedBy); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal verify achievement",
		})
	}

	// Refresh data
	achievement, _ = s.achievementRepo.GetAchievementByID(achievementID)
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil diverify",
		Data:    achievement,
	})
}

func (s *achievementServiceImpl) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	rejectedBy := c.Locals("userID").(string)

	type RejectRequest struct {
		RejectionNote string `json:"rejection_note"`
	}

	var req RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil || achievement == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi tidak ditemukan",
		})
	}

	if achievement.Status != model.AchievementStatusSubmitted {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "hanya prestasi submitted yang bisa direject",
		})
	}

	student, _ := s.studentRepo.GetStudentByID(achievement.StudentID)
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	if student.AdvisorID != rejectedBy {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "anda bukan advisor dari student ini",
		})
	}

	if err := s.achievementRepo.RejectAchievement(achievementID, rejectedBy, req.RejectionNote); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal reject achievement",
		})
	}

	// Refresh data
	achievement, _ = s.achievementRepo.GetAchievementByID(achievementID)
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil direject",
		Data:    achievement,
	})
}

func (s *achievementServiceImpl) DeleteAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	studentID := c.Locals("userID").(string)

	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil || achievement == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi tidak ditemukan",
		})
	}

	if achievement.StudentID != studentID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi bukan milik anda",
		})
	}

	if achievement.Status != model.AchievementStatusDraft {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "hanya prestasi draft yang bisa dihapus",
		})
	}

	if err := s.achievementRepo.DeleteAchievement(achievementID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal menghapus achievement",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil dihapus",
	})
}

func (s *achievementServiceImpl) GetAdviseeAchievements(c *fiber.Ctx) error {
	advisorID := c.Locals("userID").(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	students, err := s.studentRepo.GetStudentsByAdvisorID(advisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil students",
		})
	}

	if len(students) == 0 {
		return c.Status(fiber.StatusOK).JSON(model.APIResponse{
			Status:  "success",
			Message: "achievements berhasil diambil",
			Data: map[string]interface{}{
				"achievements": []model.Achievement{},
				"pagination": map[string]interface{}{
					"page":       page,
					"page_size":  pageSize,
					"total":      0,
					"total_page": 0,
				},
			},
		})
	}

	var achievements []*model.Achievement
	if status != "" {
		achievements, err = s.achievementRepo.GetAchievementsByStatus(status)
	} else {
		achievements, _, err = s.achievementRepo.GetAllAchievements(1, 1000)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil achievements",
		})
	}

	adviseeMap := make(map[string]bool)
	for _, student := range students {
		adviseeMap[student.ID] = true
	}

	var filtered []*model.Achievement
	for _, ach := range achievements {
		if adviseeMap[ach.StudentID] {
			filtered = append(filtered, ach)
		}
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		filtered = []*model.Achievement{}
	} else {
		if end > total {
			end = total
		}
		filtered = filtered[start:end]
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievements berhasil diambil",
		Data: map[string]interface{}{
			"achievements": filtered,
			"pagination": map[string]interface{}{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + pageSize - 1) / pageSize,
			},
		},
	})
}
