package service

import (
	"strconv"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService interface {
	CreateLecturer(c *fiber.Ctx) error
	GetLecturerByID(c *fiber.Ctx) error
	GetLecturerByUserID(c *fiber.Ctx) error
	GetAllLecturers(c *fiber.Ctx) error
	GetAdvisees(c *fiber.Ctx) error
	UpdateLecturer(c *fiber.Ctx) error
	DeleteLecturer(c *fiber.Ctx) error
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

func (s *lecturerServiceImpl) CreateLecturer(c *fiber.Ctx) error {
	type CreateLecturerRequest struct {
		UserID     string `json:"user_id"`
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	}

	req := new(CreateLecturerRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	if req.UserID == "" || req.LecturerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "user_id dan lecturer_id tidak boleh kosong",
		})
	}

	existing, _ := s.lecturerRepo.GetLecturerByLecturerID(req.LecturerID)
	if existing != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "lecturer ID sudah terdaftar",
		})
	}

	lecturer := &model.Lecturer{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	if err := s.lecturerRepo.CreateLecturer(lecturer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal membuat lecturer: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "lecturer berhasil dibuat",
		Data:    lecturer,
	})
}

func (s *lecturerServiceImpl) GetLecturerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	lecturer, err := s.lecturerRepo.GetLecturerByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil lecturer: " + err.Error(),
		})
	}
	if lecturer == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "lecturer tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "lecturer berhasil diambil",
		Data:    lecturer,
	})
}

func (s *lecturerServiceImpl) GetLecturerByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "user_id tidak boleh kosong",
		})
	}

	lecturer, err := s.lecturerRepo.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil lecturer: " + err.Error(),
		})
	}
	if lecturer == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "lecturer tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "lecturer berhasil diambil",
		Data:    lecturer,
	})
}

func (s *lecturerServiceImpl) GetAllLecturers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	lecturers, total, err := s.lecturerRepo.GetAllLecturers(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil lecturers: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "lecturers berhasil diambil",
		Data: map[string]interface{}{
			"lecturers": lecturers,
			"pagination": map[string]interface{}{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + pageSize - 1) / pageSize,
			},
		},
	})
}

func (s *lecturerServiceImpl) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	if lecturerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "lecturer_id tidak boleh kosong",
		})
	}

	advisees, err := s.studentRepo.GetStudentsByAdvisorID(lecturerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil advisees: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "advisees berhasil diambil",
		Data:    advisees,
	})
}

func (s *lecturerServiceImpl) UpdateLecturer(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateLecturerRequest struct {
		Department string `json:"department"`
	}

	req := new(UpdateLecturerRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	lecturer, err := s.lecturerRepo.GetLecturerByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil lecturer: " + err.Error(),
		})
	}
	if lecturer == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "lecturer tidak ditemukan",
		})
	}

	if req.Department != "" {
		lecturer.Department = req.Department
		if err := s.lecturerRepo.UpdateLecturer(lecturer); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
				Status:  "error",
				Message: "gagal mengupdate lecturer: " + err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "lecturer berhasil diupdate",
	})
}

func (s *lecturerServiceImpl) DeleteLecturer(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	lecturer, err := s.lecturerRepo.GetLecturerByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil lecturer: " + err.Error(),
		})
	}
	if lecturer == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "lecturer tidak ditemukan",
		})
	}

	if err := s.lecturerRepo.DeleteLecturer(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal menghapus lecturer: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "lecturer berhasil dihapus",
	})
}
