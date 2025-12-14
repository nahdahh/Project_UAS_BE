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
