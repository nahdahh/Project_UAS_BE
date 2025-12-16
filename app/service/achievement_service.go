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
	GetAchievementHistory(c *fiber.Ctx) error
	UploadAttachment(c *fiber.Ctx) error
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

// GetAllAchievements godoc
// @Summary Dapatkan semua prestasi
// @Description Mengambil daftar semua prestasi dengan filter dan pagination
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman" default(1)
// @Param page_size query int false "Jumlah data per halaman" default(10)
// @Param status query string false "Filter berdasarkan status" Enums(draft, submitted, verified, rejected)
// @Param achievement_type query string false "Filter berdasarkan tipe" Enums(academic, competition, organization, publication, certification, other)
// @Param student_id query string false "Filter berdasarkan student ID"
// @Param start_date query string false "Filter tanggal mulai (YYYY-MM-DD)"
// @Param end_date query string false "Filter tanggal akhir (YYYY-MM-DD)"
// @Param sort_by query string false "Field untuk sorting" default(created_at)
// @Param sort_order query string false "Urutan sorting" Enums(ASC, DESC) default(DESC)
// @Success 200 {object} model.APIResponse{data=object{achievements=[]model.AchievementWithReference,filters=object,pagination=object}} "Prestasi berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements [get]
func (s *achievementServiceImpl) GetAllAchievements(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status")
	achievementType := c.Query("achievement_type")
	studentID := c.Query("student_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "DESC")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Build filters map
	filters := make(map[string]interface{})
	if status != "" {
		filters["status"] = status
	}
	if achievementType != "" {
		filters["achievement_type"] = achievementType
	}
	if studentID != "" {
		filters["student_id"] = studentID
	}
	if startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate != "" {
		filters["end_date"] = endDate
	}

	var achievements []*model.AchievementWithReference
	var total int
	var err error

	// Use filtered query if any filters are provided
	if len(filters) > 0 || sortBy != "created_at" || sortOrder != "DESC" {
		achievements, total, err = s.achievementRepo.GetAchievementsWithFilters(page, pageSize, filters, sortBy, sortOrder)
	} else {
		achievements, total, err = s.achievementRepo.GetAllAchievements(page, pageSize)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil achievements: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievements berhasil diambil",
		Data: map[string]interface{}{
			"achievements": achievements,
			"filters": map[string]interface{}{
				"status":           status,
				"achievement_type": achievementType,
				"student_id":       studentID,
				"start_date":       startDate,
				"end_date":         endDate,
				"sort_by":          sortBy,
				"sort_order":       sortOrder,
			},
			"pagination": map[string]interface{}{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + pageSize - 1) / pageSize,
			},
		},
	})
}

// GetAchievementDetail godoc
// @Summary Dapatkan detail prestasi
// @Description Mengambil detail prestasi berdasarkan ID
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} model.APIResponse{data=model.AchievementWithReference} "Detail prestasi berhasil diambil"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 404 {object} model.APIResponse "Prestasi tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id} [get]
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

	if achievement.Status == model.AchievementStatusDeleted {
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

	if role == "Dosen Wali" {
		student, err := s.studentRepo.GetStudentByID(achievement.StudentID)
		if err != nil || student == nil || student.AdvisorID != studentID {
			return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
				Status:  "error",
				Message: "anda tidak memiliki akses ke prestasi ini",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil diambil",
		Data:    achievement,
	})
}

// CreateAchievement godoc
// @Summary Buat prestasi baru
// @Description Membuat prestasi baru (hanya mahasiswa)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body model.CreateAchievementRequest true "Data prestasi"
// @Success 201 {object} model.APIResponse{data=model.Achievement} "Prestasi berhasil dibuat"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "Student tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements [post]
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

	// Create achievement (will be saved to MongoDB + PostgreSQL reference)
	achievement := &model.Achievement{
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	result, err := s.achievementRepo.CreateAchievement(achievement, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal membuat achievement",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil dibuat",
		Data:    result,
	})
}

// UpdateAchievement godoc
// @Summary Update prestasi
// @Description Memperbarui data prestasi (hanya yang berstatus draft)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param body body model.UpdateAchievementRequest true "Data prestasi yang diupdate"
// @Success 200 {object} model.APIResponse{data=model.Achievement} "Prestasi berhasil diupdate"
// @Failure 400 {object} model.APIResponse "Format request tidak valid atau status bukan draft"
// @Failure 401 {object} model.APIResponse "Prestasi bukan milik anda"
// @Failure 403 {object} model.APIResponse "Dosen wali tidak dapat mengedit prestasi"
// @Failure 404 {object} model.APIResponse "Prestasi tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id} [put]
func (s *achievementServiceImpl) UpdateAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	studentID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

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

	if role == "Mahasiswa" && achievement.StudentID != studentID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi bukan milik anda",
		})
	}

	if role == "Dosen Wali" {
		return c.Status(fiber.StatusForbidden).JSON(model.APIResponse{
			Status:  "error",
			Message: "dosen wali tidak dapat mengedit prestasi",
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

	if err := s.achievementRepo.UpdateAchievement(achievementID, &achievement.Achievement); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal update achievement",
		})
	}

	// Refresh data
	achievement, _ = s.achievementRepo.GetAchievementByID(achievementID)
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement berhasil diupdate",
		Data:    achievement,
	})
}

// SubmitAchievement godoc
// @Summary Submit prestasi untuk verifikasi
// @Description Mengubah status prestasi dari draft menjadi submitted
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} model.APIResponse{data=model.AchievementWithReference} "Prestasi berhasil disubmit"
// @Failure 400 {object} model.APIResponse "Hanya prestasi draft yang bisa disubmit"
// @Failure 401 {object} model.APIResponse "Prestasi bukan milik anda"
// @Failure 404 {object} model.APIResponse "Prestasi tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id}/submit [post]
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

// VerifyAchievement godoc
// @Summary Verifikasi prestasi
// @Description Memverifikasi prestasi yang telah disubmit (dosen wali/admin)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} model.APIResponse{data=model.AchievementWithReference} "Prestasi berhasil diverifikasi"
// @Failure 400 {object} model.APIResponse "Hanya prestasi submitted yang bisa diverify"
// @Failure 401 {object} model.APIResponse "Anda bukan advisor dari student ini"
// @Failure 403 {object} model.APIResponse "Mahasiswa tidak dapat memverifikasi prestasi"
// @Failure 404 {object} model.APIResponse "Prestasi tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id}/verify [post]
func (s *achievementServiceImpl) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	verifiedBy := c.Locals("userID").(string)
	role := c.Locals("role").(string)

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

	if role == "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(model.APIResponse{
			Status:  "error",
			Message: "mahasiswa tidak dapat memverifikasi prestasi",
		})
	}

	if role == "Dosen Wali" {
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

// RejectAchievement godoc
// @Summary Tolak prestasi
// @Description Menolak prestasi yang telah disubmit dengan catatan (dosen wali/admin)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param body body object{rejection_note=string} true "Catatan penolakan"
// @Success 200 {object} model.APIResponse{data=model.AchievementWithReference} "Prestasi berhasil ditolak"
// @Failure 400 {object} model.APIResponse "Format request tidak valid atau hanya prestasi submitted yang bisa direject"
// @Failure 401 {object} model.APIResponse "Anda bukan advisor dari student ini"
// @Failure 403 {object} model.APIResponse "Mahasiswa tidak dapat menolak prestasi"
// @Failure 404 {object} model.APIResponse "Prestasi tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id}/reject [post]
func (s *achievementServiceImpl) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	rejectedBy := c.Locals("userID").(string)
	role := c.Locals("role").(string)

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

	if role == "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(model.APIResponse{
			Status:  "error",
			Message: "mahasiswa tidak dapat menolak prestasi",
		})
	}

	if role == "Dosen Wali" {
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

// DeleteAchievement godoc
// @Summary Hapus prestasi
// @Description Menghapus prestasi (hanya yang berstatus draft)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} model.APIResponse "Prestasi berhasil dihapus"
// @Failure 400 {object} model.APIResponse "Hanya prestasi draft yang bisa dihapus"
// @Failure 401 {object} model.APIResponse "Prestasi bukan milik anda"
// @Failure 403 {object} model.APIResponse "Dosen wali tidak dapat menghapus prestasi"
// @Failure 404 {object} model.APIResponse "Prestasi tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id} [delete]
func (s *achievementServiceImpl) DeleteAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	studentID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

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

	if role == "Mahasiswa" && achievement.StudentID != studentID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "prestasi bukan milik anda",
		})
	}

	if role == "Dosen Wali" {
		return c.Status(fiber.StatusForbidden).JSON(model.APIResponse{
			Status:  "error",
			Message: "dosen wali tidak dapat menghapus prestasi",
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

// GetAdviseeAchievements godoc
// @Summary Dapatkan prestasi anak bimbingan
// @Description Mengambil daftar prestasi mahasiswa yang dibimbing (khusus dosen wali)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman" default(1)
// @Param page_size query int false "Jumlah data per halaman" default(10)
// @Param status query string false "Filter berdasarkan status"
// @Success 200 {object} model.APIResponse{data=object{achievements=[]model.AchievementWithReference,pagination=object}} "Prestasi anak bimbingan berhasil diambil"
// @Failure 403 {object} model.APIResponse "Mahasiswa tidak dapat mengakses endpoint ini"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/advisee/list [get]
func (s *achievementServiceImpl) GetAdviseeAchievements(c *fiber.Ctx) error {
	advisorID := c.Locals("userID").(string)
	role := c.Locals("role").(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	if role == "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(model.APIResponse{
			Status:  "error",
			Message: "mahasiswa tidak dapat mengakses endpoint ini",
		})
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
				"achievements": []*model.AchievementWithReference{},
				"pagination": map[string]interface{}{
					"page":       page,
					"page_size":  pageSize,
					"total":      0,
					"total_page": 0,
				},
			},
		})
	}

	var achievements []*model.AchievementWithReference
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

	var filtered []*model.AchievementWithReference
	for _, ach := range achievements {
		if adviseeMap[ach.StudentID] {
			filtered = append(filtered, ach)
		}
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		filtered = []*model.AchievementWithReference{}
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

// GetAchievementHistory godoc
// @Summary Dapatkan riwayat prestasi
// @Description Mengambil riwayat perubahan status prestasi
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} model.APIResponse{data=[]model.AchievementHistory} "Riwayat berhasil diambil"
// @Failure 400 {object} model.APIResponse "Achievement ID tidak boleh kosong"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 404 {object} model.APIResponse "Achievement tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id}/history [get]
func (s *achievementServiceImpl) GetAchievementHistory(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("userID").(string)

	if achievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "achievement_id tidak boleh kosong",
		})
	}

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
			Message: "achievement tidak ditemukan",
		})
	}

	role := c.Locals("role").(string)
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "unauthorized",
		})
	}

	history, err := s.achievementRepo.GetAchievementHistory(achievementID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil history",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "history berhasil diambil",
		Data:    history,
	})
}

// UploadAttachment godoc
// @Summary Upload lampiran prestasi
// @Description Mengupload file lampiran untuk prestasi
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param file formData file true "File lampiran"
// @Success 201 {object} model.APIResponse{data=model.AchievementAttachment} "Lampiran berhasil diupload"
// @Failure 400 {object} model.APIResponse "File harus diupload"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 404 {object} model.APIResponse "Achievement tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /achievements/{id}/attachments [post]
func (s *achievementServiceImpl) UploadAttachment(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("userID").(string)

	if achievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "achievement_id tidak boleh kosong",
		})
	}

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
			Message: "achievement tidak ditemukan",
		})
	}

	if achievement.StudentID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(model.APIResponse{
			Status:  "error",
			Message: "unauthorized",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "file harus diupload",
		})
	}

	// Save file to storage (simplified - in production use cloud storage)
	fileName := uuid.New().String() + "_" + file.Filename
	filePath := "./uploads/" + fileName

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal menyimpan file",
		})
	}

	attachment := &model.AchievementAttachment{
		ID:            uuid.New().String(),
		AchievementID: achievementID,
		FileName:      file.Filename,
		FileURL:       "/uploads/" + fileName,
		FileSize:      file.Size,
		FileType:      file.Header.Get("Content-Type"),
		UploadedBy:    userID,
		CreatedAt:     time.Now(),
	}

	if err := s.achievementRepo.CreateAttachment(attachment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal menyimpan attachment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "attachment berhasil diupload",
		Data:    attachment,
	})
}
