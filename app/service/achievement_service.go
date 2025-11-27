package service

import (
	"errors"
	"strconv"
	"time"
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementService interface {
	CreateAchievementReference(studentID string, achievementType string) (*model.AchievementReference, error)
	SubmitAchievementForVerification(referenceID, studentID string) error
	VerifyAchievement(referenceID, advisorID string) error
	RejectAchievement(referenceID, advisorID, rejectionNote string) error
	GetStudentAchievements(studentID string, page, pageSize int) ([]*model.AchievementReference, int, error)
	GetAchievementDetail(referenceID, studentID string) (*model.AchievementReference, error)
	GetAdviseeAchievements(advisorID, status string, page, pageSize int) ([]*model.AchievementReference, int, error)
	GetAllAchievements(page, pageSize int) ([]*model.AchievementReference, int, error)
	DeleteAchievementReference(referenceID, studentID string) error
	// HTTP Handler methods
	HandleGetAll(c *fiber.Ctx) error
	HandleGetByID(c *fiber.Ctx) error
	HandleCreate(c *fiber.Ctx) error
	HandleSubmit(c *fiber.Ctx) error
	HandleVerify(c *fiber.Ctx) error
	HandleReject(c *fiber.Ctx) error
	HandleDelete(c *fiber.Ctx) error
	HandleGetAdviseeAchievements(c *fiber.Ctx) error
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

// CreateAchievementReference membuat referensi prestasi baru dengan status draft
func (s *achievementServiceImpl) CreateAchievementReference(studentID string, achievementType string) (*model.AchievementReference, error) {
	// Validasi: student harus ada di database
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("student tidak ditemukan")
	}

	// Business logic: buat achievement baru dengan status draft
	ref := &model.AchievementReference{
		ID:        uuid.New().String(),
		StudentID: studentID,
		Status:    model.AchievementStatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Simpan ke database via repository
	err = s.achievementRepo.CreateAchievement(ref)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

// SubmitAchievementForVerification mengubah status prestasi dari draft menjadi submitted
func (s *achievementServiceImpl) SubmitAchievementForVerification(referenceID, studentID string) error {
	// Ambil achievement dari database
	ref, err := s.achievementRepo.GetAchievementByID(referenceID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("prestasi tidak ditemukan")
	}

	// Business rule: hanya pemilik yang bisa submit
	if ref.StudentID != studentID {
		return errors.New("prestasi bukan milik anda")
	}

	// Business rule: hanya draft yang bisa disubmit
	if ref.Status != model.AchievementStatusDraft {
		return errors.New("hanya prestasi draft yang bisa disubmit")
	}

	// Update status ke submitted
	now := time.Now()
	ref.Status = model.AchievementStatusSubmitted
	ref.SubmittedAt = &now
	ref.UpdatedAt = now

	return s.achievementRepo.UpdateAchievement(ref)
}

// VerifyAchievement mengubah status prestasi dari submitted menjadi verified
func (s *achievementServiceImpl) VerifyAchievement(referenceID, advisorID string) error {
	// Ambil achievement dari database
	ref, err := s.achievementRepo.GetAchievementByID(referenceID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("prestasi tidak ditemukan")
	}

	// Business rule: hanya submitted yang bisa diverify
	if ref.Status != model.AchievementStatusSubmitted {
		return errors.New("hanya prestasi submitted yang bisa diverify")
	}

	// Business rule: hanya advisor dari student yang bisa verify
	student, err := s.studentRepo.GetStudentByID(ref.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student tidak ditemukan")
	}
	if student.AdvisorID != advisorID {
		return errors.New("anda bukan advisor dari student ini")
	}

	// Update status ke verified
	now := time.Now()
	ref.Status = model.AchievementStatusVerified
	ref.VerifiedBy = &advisorID
	ref.VerifiedAt = &now
	ref.UpdatedAt = now

	return s.achievementRepo.UpdateAchievement(ref)
}

// RejectAchievement menolak prestasi dengan alasan penolakan
func (s *achievementServiceImpl) RejectAchievement(referenceID, advisorID, rejectionNote string) error {
	// Ambil achievement dari database
	ref, err := s.achievementRepo.GetAchievementByID(referenceID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("prestasi tidak ditemukan")
	}

	// Business rule: hanya submitted yang bisa direject
	if ref.Status != model.AchievementStatusSubmitted {
		return errors.New("hanya prestasi submitted yang bisa direject")
	}

	// Business rule: hanya advisor dari student yang bisa reject
	student, err := s.studentRepo.GetStudentByID(ref.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student tidak ditemukan")
	}
	if student.AdvisorID != advisorID {
		return errors.New("anda bukan advisor dari student ini")
	}

	// Update status ke rejected
	now := time.Now()
	ref.Status = model.AchievementStatusRejected
	ref.VerifiedBy = &advisorID
	ref.RejectionNote = &rejectionNote
	ref.UpdatedAt = now

	return s.achievementRepo.UpdateAchievement(ref)
}

// GetStudentAchievements mengambil semua prestasi milik student dengan pagination
func (s *achievementServiceImpl) GetStudentAchievements(studentID string, page, pageSize int) ([]*model.AchievementReference, int, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	refs, err := s.achievementRepo.GetAchievementsByStudentID(studentID)
	if err != nil {
		return nil, 0, err
	}

	return refs, len(refs), nil
}

// GetAchievementDetail mengambil detail prestasi dengan cek ownership
func (s *achievementServiceImpl) GetAchievementDetail(referenceID, studentID string) (*model.AchievementReference, error) {
	ref, err := s.achievementRepo.GetAchievementByID(referenceID)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, errors.New("prestasi tidak ditemukan")
	}

	// Business rule: hanya pemilik yang bisa lihat detail
	if ref.StudentID != studentID {
		return nil, errors.New("unauthorized")
	}

	return ref, nil
}

// GetAdviseeAchievements mengambil prestasi dari semua mahasiswa bimbingan advisor
func (s *achievementServiceImpl) GetAdviseeAchievements(advisorID, status string, page, pageSize int) ([]*model.AchievementReference, int, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Ambil semua student dari advisor
	students, err := s.studentRepo.GetStudentsByAdvisorID(advisorID)
	if err != nil {
		return nil, 0, err
	}

	if len(students) == 0 {
		return []*model.AchievementReference{}, 0, nil
	}

	// Ambil prestasi berdasarkan status jika ada
	var refs []*model.AchievementReference
	if status != "" {
		refs, err = s.achievementRepo.GetAchievementsByStatus(status)
	} else {
		// Jika tidak ada filter status, ambil semua
		_, _, err = s.GetAllAchievements(1, 1000)
		if err != nil {
			return nil, 0, err
		}
	}
	if err != nil {
		return nil, 0, err
	}

	// Filter hanya achievement dari student-student yang di-advise oleh advisor
	adviseeMap := make(map[string]bool)
	for _, st := range students {
		adviseeMap[st.ID] = true
	}

	var filtered []*model.AchievementReference
	for _, ref := range refs {
		if adviseeMap[ref.StudentID] {
			filtered = append(filtered, ref)
		}
	}

	// Apply pagination
	total := len(filtered)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []*model.AchievementReference{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// GetAllAchievements mengambil semua prestasi di sistem (admin only)
func (s *achievementServiceImpl) GetAllAchievements(page, pageSize int) ([]*model.AchievementReference, int, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	refs, total, err := s.achievementRepo.GetAllAchievements(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return refs, total, nil
}

// DeleteAchievementReference menghapus prestasi (soft delete) jika masih draft
func (s *achievementServiceImpl) DeleteAchievementReference(referenceID, studentID string) error {
	// Ambil achievement dari database
	ref, err := s.achievementRepo.GetAchievementByID(referenceID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("prestasi tidak ditemukan")
	}

	// Business rule: hanya pemilik yang bisa hapus
	if ref.StudentID != studentID {
		return errors.New("prestasi bukan milik anda")
	}

	// Business rule: hanya draft yang bisa dihapus
	if ref.Status != model.AchievementStatusDraft {
		return errors.New("hanya prestasi draft yang bisa dihapus")
	}

	return s.achievementRepo.DeleteAchievement(referenceID)
}

// ===== HTTP HANDLER METHODS (Route Layer) =====
// Handlers hanya parse input → panggil business logic → return response

// HandleGetAll menangani GET /api/v1/achievements
func (s *achievementServiceImpl) HandleGetAll(c *fiber.Ctx) error {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	achievements, total, err := s.GetAllAchievements(page, pageSize)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data achievement berhasil diambil", fiber.Map{
		"data":  achievements,
		"total": total,
		"page":  page,
	})
}

// HandleGetByID menangani GET /api/v1/achievements/:id
func (s *achievementServiceImpl) HandleGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	studentID := c.Locals("user_id").(string)

	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	achievement, err := s.GetAchievementDetail(id, studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "achievement tidak ditemukan: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data achievement berhasil diambil", achievement)
}

// HandleCreate menangani POST /api/v1/achievements
func (s *achievementServiceImpl) HandleCreate(c *fiber.Ctx) error {
	studentID := c.Locals("user_id").(string)

	var req struct {
		AchievementType string `json:"achievement_type"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
	}

	achievement, err := s.CreateAchievementReference(studentID, req.AchievementType)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal membuat achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "achievement berhasil dibuat", achievement)
}

// HandleSubmit menangani POST /api/v1/achievements/:id/submit
func (s *achievementServiceImpl) HandleSubmit(c *fiber.Ctx) error {
	id := c.Params("id")
	studentID := c.Locals("user_id").(string)

	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	err := s.SubmitAchievementForVerification(id, studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal submit achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "achievement berhasil disubmit", nil)
}

// HandleVerify menangani POST /api/v1/achievements/:id/verify
func (s *achievementServiceImpl) HandleVerify(c *fiber.Ctx) error {
	id := c.Params("id")
	advisorID := c.Locals("user_id").(string)

	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	err := s.VerifyAchievement(id, advisorID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal verify achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "achievement berhasil diverifikasi", nil)
}

// HandleReject menangani POST /api/v1/achievements/:id/reject
func (s *achievementServiceImpl) HandleReject(c *fiber.Ctx) error {
	id := c.Params("id")
	advisorID := c.Locals("user_id").(string)

	var req struct {
		RejectionNote string `json:"rejection_note"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
	}

	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	err := s.RejectAchievement(id, advisorID, req.RejectionNote)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal reject achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "achievement berhasil ditolak", nil)
}

// HandleDelete menangani DELETE /api/v1/achievements/:id
func (s *achievementServiceImpl) HandleDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	studentID := c.Locals("user_id").(string)

	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	err := s.DeleteAchievementReference(id, studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal menghapus achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "achievement berhasil dihapus", nil)
}

// HandleGetAdviseeAchievements menangani GET /api/v1/achievements/advisee/list
func (s *achievementServiceImpl) HandleGetAdviseeAchievements(c *fiber.Ctx) error {
	advisorID := c.Locals("user_id").(string)
	status := c.Query("status", "")

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	achievements, total, err := s.GetAdviseeAchievements(advisorID, status, page, pageSize)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data achievement: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data achievement advisee berhasil diambil", fiber.Map{
		"data":  achievements,
		"total": total,
		"page":  page,
	})
}
