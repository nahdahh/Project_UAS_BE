package service

import (
	"strconv"
	"strings"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	CreateStudent(c *fiber.Ctx) error
	GetStudentByID(c *fiber.Ctx) error
	GetStudentByUserID(c *fiber.Ctx) error
	GetAllStudents(c *fiber.Ctx) error
	GetStudentsByAdvisor(c *fiber.Ctx) error
	UpdateStudent(c *fiber.Ctx) error
	DeleteStudent(c *fiber.Ctx) error
	GetStudentAchievements(c *fiber.Ctx) error
	UpdateAdvisor(c *fiber.Ctx) error
}

type studentServiceImpl struct {
	studentRepo     repository.StudentRepository
	achievementRepo repository.AchievementRepository
}

func NewStudentService(studentRepo repository.StudentRepository, achievementRepo repository.AchievementRepository) StudentService {
	return &studentServiceImpl{
		studentRepo:     studentRepo,
		achievementRepo: achievementRepo,
	}
}

// CreateStudent godoc
// @Summary Buat data mahasiswa baru
// @Description Membuat data mahasiswa baru dengan user ID yang valid
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{user_id=string,student_id=string,program_study=string,academic_year=string,advisor_id=string} true "Data mahasiswa"
// @Success 201 {object} model.APIResponse{data=model.Student} "Mahasiswa berhasil dibuat"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students [post]
func (s *studentServiceImpl) CreateStudent(c *fiber.Ctx) error {
	type CreateStudentRequest struct {
		UserID       string `json:"user_id"`
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
		AdvisorID    string `json:"advisor_id"`
	}

	req := new(CreateStudentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	// Trim spaces
	req.UserID = strings.TrimSpace(req.UserID)
	req.StudentID = strings.TrimSpace(req.StudentID)
	req.AdvisorID = strings.TrimSpace(req.AdvisorID)

	if req.UserID == "" || req.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "user_id dan student_id tidak boleh kosong",
		})
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "user_id harus berformat UUID yang valid",
		})
	}

	if req.AdvisorID != "" {
		if _, err := uuid.Parse(req.AdvisorID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
				Status:  "error",
				Message: "advisor_id harus berformat UUID yang valid",
			})
		}
	}

	existing, _ := s.studentRepo.GetStudentByStudentID(req.StudentID)
	if existing != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "student ID (NIM) sudah terdaftar",
		})
	}

	student := &model.Student{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	}

	if err := s.studentRepo.CreateStudent(student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal membuat student: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "student berhasil dibuat",
		Data:    student,
	})
}

// GetStudentByID godoc
// @Summary Dapatkan mahasiswa berdasarkan ID
// @Description Mengambil data mahasiswa berdasarkan ID
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} model.APIResponse{data=model.Student} "Data mahasiswa berhasil diambil"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/{id} [get]
func (s *studentServiceImpl) GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student: " + err.Error(),
		})
	}
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "student berhasil diambil",
		Data:    student,
	})
}

// GetStudentByUserID godoc
// @Summary Dapatkan mahasiswa berdasarkan User ID
// @Description Mengambil data mahasiswa berdasarkan User ID
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} model.APIResponse{data=model.Student} "Data mahasiswa berhasil diambil"
// @Failure 400 {object} model.APIResponse "User ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/user/{user_id} [get]
func (s *studentServiceImpl) GetStudentByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "user_id tidak boleh kosong",
		})
	}

	student, err := s.studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student: " + err.Error(),
		})
	}
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "student berhasil diambil",
		Data:    student,
	})
}

// GetAllStudents godoc
// @Summary Dapatkan semua mahasiswa
// @Description Mengambil daftar semua mahasiswa dengan pagination
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman" default(1)
// @Param page_size query int false "Jumlah data per halaman" default(10)
// @Success 200 {object} model.APIResponse{data=object{students=[]model.Student,pagination=object}} "Daftar mahasiswa berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students [get]
func (s *studentServiceImpl) GetAllStudents(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	students, total, err := s.studentRepo.GetAllStudents(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil students: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "students berhasil diambil",
		Data: map[string]interface{}{
			"students": students,
			"pagination": map[string]interface{}{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + pageSize - 1) / pageSize,
			},
		},
	})
}

// GetStudentsByAdvisor godoc
// @Summary Dapatkan mahasiswa berdasarkan dosen wali
// @Description Mengambil daftar mahasiswa yang dibimbing oleh dosen wali tertentu
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param advisor_id path string true "Advisor ID (Lecturer ID)"
// @Success 200 {object} model.APIResponse{data=[]model.Student} "Daftar mahasiswa berhasil diambil"
// @Failure 400 {object} model.APIResponse "Advisor ID tidak boleh kosong"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/advisor/{advisor_id} [get]
func (s *studentServiceImpl) GetStudentsByAdvisor(c *fiber.Ctx) error {
	advisorID := c.Params("advisor_id")

	if advisorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "advisor_id tidak boleh kosong",
		})
	}

	students, err := s.studentRepo.GetStudentsByAdvisorID(advisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil students: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "students berhasil diambil",
		Data:    students,
	})
}

// UpdateStudent godoc
// @Summary Update data mahasiswa
// @Description Memperbarui data mahasiswa
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param body body object{program_study=string,academic_year=string,advisor_id=string} true "Data yang diupdate"
// @Success 200 {object} model.APIResponse "Mahasiswa berhasil diupdate"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/{id} [put]
func (s *studentServiceImpl) UpdateStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateStudentRequest struct {
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
		AdvisorID    string `json:"advisor_id"`
	}

	req := new(UpdateStudentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	// Trim spaces
	req.AdvisorID = strings.TrimSpace(req.AdvisorID)

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student: " + err.Error(),
		})
	}
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	updateData := &model.Student{ID: id}
	if req.ProgramStudy != "" {
		updateData.ProgramStudy = req.ProgramStudy
	}
	if req.AcademicYear != "" {
		updateData.AcademicYear = req.AcademicYear
	}
	if req.AdvisorID != "" {
		if _, err := uuid.Parse(req.AdvisorID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
				Status:  "error",
				Message: "advisor_id harus berformat UUID yang valid",
			})
		}
		updateData.AdvisorID = req.AdvisorID
	}

	if err := s.studentRepo.UpdateStudent(updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengupdate student: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "student berhasil diupdate",
	})
}

// DeleteStudent godoc
// @Summary Hapus data mahasiswa
// @Description Menghapus data mahasiswa berdasarkan ID
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} model.APIResponse "Mahasiswa berhasil dihapus"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/{id} [delete]
func (s *studentServiceImpl) DeleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student: " + err.Error(),
		})
	}
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	if err := s.studentRepo.DeleteStudent(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal menghapus student: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "student berhasil dihapus",
	})
}

// GetStudentAchievements godoc
// @Summary Dapatkan prestasi mahasiswa
// @Description Mengambil daftar prestasi mahasiswa berdasarkan ID
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} model.APIResponse{data=[]model.AchievementWithReference} "Daftar prestasi berhasil diambil"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/{id}/achievements [get]
func (s *studentServiceImpl) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")

	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student: " + err.Error(),
		})
	}
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	achievements, err := s.achievementRepo.GetAchievementsByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil achievements: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievements berhasil diambil",
		Data:    achievements,
	})
}

// UpdateAdvisor godoc
// @Summary Update dosen wali mahasiswa
// @Description Mengubah dosen wali mahasiswa
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param body body object{advisor_id=string} true "Advisor ID baru"
// @Success 200 {object} model.APIResponse "Dosen wali berhasil diupdate"
// @Failure 400 {object} model.APIResponse "Student ID dan advisor ID harus diisi"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /students/{id}/advisor [put]
func (s *studentServiceImpl) UpdateAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	type UpdateAdvisorRequest struct {
		AdvisorID string `json:"advisor_id"`
	}

	req := new(UpdateAdvisorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	if studentID == "" || req.AdvisorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "student_id dan advisor_id harus diisi",
		})
	}

	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student: " + err.Error(),
		})
	}
	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "student tidak ditemukan",
		})
	}

	if err := s.studentRepo.UpdateAdvisor(studentID, req.AdvisorID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal update advisor: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "advisor berhasil diupdate",
	})
}
