package service

import (
	"errors"
	"strconv"
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService interface {
	CreateLecturer(userID, lecturerID, department string) (*model.Lecturer, error)
	GetLecturerByID(id string) (*model.Lecturer, error)
	GetLecturerByUserID(userID string) (*model.Lecturer, error)
	GetAllLecturers(page, pageSize int) ([]*model.LecturerWithUser, int, error)
	GetAdvisees(lecturerID string) ([]*model.StudentWithUser, error)
	UpdateLecturer(id, department string) error
	DeleteLecturer(id string) error
	// HTTP Handler methods
	HandleGetAll(c *fiber.Ctx) error
	HandleGetByID(c *fiber.Ctx) error
	HandleGetAdvisees(c *fiber.Ctx) error
	HandleCreate(c *fiber.Ctx) error
	HandleUpdate(c *fiber.Ctx) error
	HandleDelete(c *fiber.Ctx) error
}

type lecturerServiceImpl struct {
	lecturerRepo repository.LecturerRepository
	studentRepo  repository.StudentRepository
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
) LecturerService {
	return &lecturerServiceImpl{
		lecturerRepo: lecturerRepo,
		studentRepo:  studentRepo,
	}
}

// CreateLecturer membuat dosen baru
func (s *lecturerServiceImpl) CreateLecturer(userID, lecturerID, department string) (*model.Lecturer, error) {
	// Validasi input wajib
	if userID == "" || lecturerID == "" {
		return nil, errors.New("user_id dan lecturer_id tidak boleh kosong")
	}

	// Business rule: lecturer_id harus unik
	existing, _ := s.lecturerRepo.GetLecturerByLecturerID(lecturerID)
	if existing != nil {
		return nil, errors.New("lecturer ID sudah terdaftar")
	}

	lecturer := &model.Lecturer{
		ID:         uuid.New().String(),
		UserID:     userID,
		LecturerID: lecturerID,
		Department: department,
	}

	err := s.lecturerRepo.CreateLecturer(lecturer)
	if err != nil {
		return nil, err
	}

	return lecturer, nil
}

// GetLecturerByID mengambil detail dosen berdasarkan ID
func (s *lecturerServiceImpl) GetLecturerByID(id string) (*model.Lecturer, error) {
	lecturer, err := s.lecturerRepo.GetLecturerByID(id)
	if err != nil {
		return nil, err
	}
	if lecturer == nil {
		return nil, errors.New("lecturer tidak ditemukan")
	}

	return lecturer, nil
}

// GetLecturerByUserID mengambil data dosen berdasarkan user ID
func (s *lecturerServiceImpl) GetLecturerByUserID(userID string) (*model.Lecturer, error) {
	lecturer, err := s.lecturerRepo.GetLecturerByUserID(userID)
	if err != nil {
		return nil, err
	}
	if lecturer == nil {
		return nil, errors.New("lecturer tidak ditemukan")
	}

	return lecturer, nil
}

// GetAllLecturers mengambil semua dosen dengan pagination
func (s *lecturerServiceImpl) GetAllLecturers(page, pageSize int) ([]*model.LecturerWithUser, int, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	lecturers, total, err := s.lecturerRepo.GetAllLecturers(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return lecturers, total, nil
}

// GetAdvisees mengambil semua mahasiswa yang dibimbing oleh dosen
func (s *lecturerServiceImpl) GetAdvisees(lecturerID string) ([]*model.StudentWithUser, error) {
	advisees, err := s.studentRepo.GetStudentsByAdvisorID(lecturerID)
	if err != nil {
		return nil, err
	}

	return advisees, nil
}

// UpdateLecturer mengupdate data dosen
func (s *lecturerServiceImpl) UpdateLecturer(id, department string) error {
	lecturer, err := s.lecturerRepo.GetLecturerByID(id)
	if err != nil {
		return err
	}
	if lecturer == nil {
		return errors.New("lecturer tidak ditemukan")
	}

	// Hanya update field yang diberikan
	if department != "" {
		lecturer.Department = department
	}

	return s.lecturerRepo.UpdateLecturer(lecturer)
}

// DeleteLecturer menghapus dosen
func (s *lecturerServiceImpl) DeleteLecturer(id string) error {
	lecturer, err := s.lecturerRepo.GetLecturerByID(id)
	if err != nil {
		return err
	}
	if lecturer == nil {
		return errors.New("lecturer tidak ditemukan")
	}

	return s.lecturerRepo.DeleteLecturer(id)
}

// ===== HTTP HANDLER METHODS (Route Layer) =====
// Handlers hanya parse input → panggil business logic → return response

// HandleGetAll menangani GET /api/v1/lecturers
func (s *lecturerServiceImpl) HandleGetAll(c *fiber.Ctx) error {
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

	lecturers, total, err := s.GetAllLecturers(page, pageSize)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data lecturer: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data lecturer berhasil diambil", fiber.Map{
		"data":  lecturers,
		"total": total,
		"page":  page,
	})
}

// HandleGetByID menangani GET /api/v1/lecturers/:id
func (s *lecturerServiceImpl) HandleGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	lecturer, err := s.GetLecturerByID(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "lecturer tidak ditemukan: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data lecturer berhasil diambil", lecturer)
}

// HandleGetAdvisees menangani GET /api/v1/lecturers/:id/advisees
func (s *lecturerServiceImpl) HandleGetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")
	if lecturerID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	advisees, err := s.GetAdvisees(lecturerID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data advisee: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data advisee berhasil diambil", advisees)
}

// HandleCreate menangani POST /api/v1/lecturers
func (s *lecturerServiceImpl) HandleCreate(c *fiber.Ctx) error {
	var req struct {
		UserID     string `json:"user_id"`
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
	}

	lecturer, err := s.CreateLecturer(req.UserID, req.LecturerID, req.Department)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal membuat lecturer: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "lecturer berhasil dibuat", lecturer)
}

// HandleUpdate menangani PUT /api/v1/lecturers/:id
func (s *lecturerServiceImpl) HandleUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	var req struct {
		Department string `json:"department"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
	}

	err := s.UpdateLecturer(id, req.Department)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal mengupdate lecturer: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "lecturer berhasil diupdate", nil)
}

// HandleDelete menangani DELETE /api/v1/lecturers/:id
func (s *lecturerServiceImpl) HandleDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	err := s.DeleteLecturer(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal menghapus lecturer: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "lecturer berhasil dihapus", nil)
}
