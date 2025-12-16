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

// CreateLecturer godoc
// @Summary Buat data dosen baru
// @Description Membuat data dosen baru dengan user ID yang valid
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{user_id=string,lecturer_id=string,department=string} true "Data dosen"
// @Success 201 {object} model.APIResponse{data=model.Lecturer} "Dosen berhasil dibuat"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers [post]
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

// GetLecturerByID godoc
// @Summary Dapatkan dosen berdasarkan ID
// @Description Mengambil data dosen berdasarkan ID
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.APIResponse{data=model.Lecturer} "Data dosen berhasil diambil"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Dosen tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers/{id} [get]
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

// GetLecturerByUserID godoc
// @Summary Dapatkan dosen berdasarkan User ID
// @Description Mengambil data dosen berdasarkan User ID
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} model.APIResponse{data=model.Lecturer} "Data dosen berhasil diambil"
// @Failure 400 {object} model.APIResponse "User ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Dosen tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers/user/{user_id} [get]
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

// GetAllLecturers godoc
// @Summary Dapatkan semua dosen
// @Description Mengambil daftar semua dosen dengan pagination
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman" default(1)
// @Param page_size query int false "Jumlah data per halaman" default(10)
// @Success 200 {object} model.APIResponse{data=object{lecturers=[]model.Lecturer,pagination=object}} "Daftar dosen berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers [get]
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

// GetAdvisees godoc
// @Summary Dapatkan anak bimbingan dosen
// @Description Mengambil daftar mahasiswa yang dibimbing oleh dosen
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.APIResponse{data=[]model.Student} "Daftar anak bimbingan berhasil diambil"
// @Failure 400 {object} model.APIResponse "Lecturer ID tidak boleh kosong"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers/{id}/advisees [get]
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

// UpdateLecturer godoc
// @Summary Update data dosen
// @Description Memperbarui data dosen
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Param body body object{department=string} true "Data yang diupdate"
// @Success 200 {object} model.APIResponse "Dosen berhasil diupdate"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "Dosen tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers/{id} [put]
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

// DeleteLecturer godoc
// @Summary Hapus data dosen
// @Description Menghapus data dosen berdasarkan ID
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.APIResponse "Dosen berhasil dihapus"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Dosen tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /lecturers/{id} [delete]
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
