package service

import (
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// StudentService interface - hanya business logic, tanpa HTTP handling
type StudentService interface {
	CreateStudent(req model.StudentCreateRequest) (map[string]interface{}, *helper.ApiError)
	GetStudentByID(id string) (map[string]interface{}, *helper.ApiError)
	GetStudentByUserID(userID string) (*model.Student, *helper.ApiError)
	GetAllStudents(page, pageSize int) (map[string]interface{}, *helper.ApiError)
	GetStudentsByAdvisor(advisorID string) ([]*model.StudentWithUser, *helper.ApiError)
	UpdateStudent(id string, req model.StudentUpdateRequest) *helper.ApiError
	DeleteStudent(id string) *helper.ApiError
	HandleGetAll(c *fiber.Ctx) error
	HandleGetByID(c *fiber.Ctx) error
	HandleCreate(c *fiber.Ctx) error
	HandleUpdate(c *fiber.Ctx) error
	HandleDelete(c *fiber.Ctx) error
}

type studentServiceImpl struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentServiceImpl{
		studentRepo: studentRepo,
	}
}

// CreateStudent membuat student baru - hanya business logic
func (s *studentServiceImpl) CreateStudent(req model.StudentCreateRequest) (map[string]interface{}, *helper.ApiError) {
	// Validasi input wajib
	if req.UserID == "" || req.StudentID == "" {
		return nil, helper.NewApiError(400, "user_id dan student_id tidak boleh kosong")
	}

	// Cek apakah student ID (NIM) sudah terdaftar
	existing, _ := s.studentRepo.GetStudentByStudentID(req.StudentID)
	if existing != nil {
		return nil, helper.NewApiError(400, "student ID (NIM) sudah terdaftar")
	}

	// Buat student baru dengan UUID
	student := &model.Student{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	}

	// Simpan ke database via repository
	if err := s.studentRepo.CreateStudent(student); err != nil {
		return nil, helper.NewApiError(500, "gagal membuat student: "+err.Error())
	}

	// Return response dalam map
	return map[string]interface{}{
		"id":            student.ID,
		"user_id":       student.UserID,
		"student_id":    student.StudentID,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
		"advisor_id":    student.AdvisorID,
		"created_at":    student.CreatedAt,
	}, nil
}

// GetStudentByID mengambil student by ID
func (s *studentServiceImpl) GetStudentByID(id string) (map[string]interface{}, *helper.ApiError) {
	if id == "" {
		return nil, helper.NewApiError(400, "id tidak boleh kosong")
	}

	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return nil, helper.NewApiError(500, "gagal mengambil student: "+err.Error())
	}
	if student == nil {
		return nil, helper.NewApiError(404, "student tidak ditemukan")
	}

	// Return response dalam map
	return map[string]interface{}{
		"id":            student.ID,
		"user_id":       student.UserID,
		"student_id":    student.StudentID,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
		"advisor_id":    student.AdvisorID,
		"created_at":    student.CreatedAt,
	}, nil
}

// GetStudentByUserID mengambil student by user ID
func (s *studentServiceImpl) GetStudentByUserID(userID string) (*model.Student, *helper.ApiError) {
	if userID == "" {
		return nil, helper.NewApiError(400, "user_id tidak boleh kosong")
	}

	student, err := s.studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return nil, helper.NewApiError(500, "gagal mengambil student: "+err.Error())
	}
	if student == nil {
		return nil, helper.NewApiError(404, "student tidak ditemukan")
	}

	return student, nil
}

// GetAllStudents mengambil semua student dengan pagination
func (s *studentServiceImpl) GetAllStudents(page, pageSize int) (map[string]interface{}, *helper.ApiError) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	students, total, err := s.studentRepo.GetAllStudents(page, pageSize)
	if err != nil {
		return nil, helper.NewApiError(500, "gagal mengambil students: "+err.Error())
	}

	// Return response dalam map dengan pagination info
	return map[string]interface{}{
		"data":      students,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	}, nil
}

// GetStudentsByAdvisor mengambil students berdasarkan advisor
func (s *studentServiceImpl) GetStudentsByAdvisor(advisorID string) ([]*model.StudentWithUser, *helper.ApiError) {
	if advisorID == "" {
		return nil, helper.NewApiError(400, "advisor_id tidak boleh kosong")
	}

	students, err := s.studentRepo.GetStudentsByAdvisorID(advisorID)
	if err != nil {
		return nil, helper.NewApiError(500, "gagal mengambil students: "+err.Error())
	}

	return students, nil
}

// UpdateStudent mengupdate student
func (s *studentServiceImpl) UpdateStudent(id string, req model.StudentUpdateRequest) *helper.ApiError {
	if id == "" {
		return helper.NewApiError(400, "id tidak boleh kosong")
	}

	// Cek student ada atau tidak
	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return helper.NewApiError(500, "gagal mengambil student: "+err.Error())
	}
	if student == nil {
		return helper.NewApiError(404, "student tidak ditemukan")
	}

	// Update field yang diberikan (optional fields)
	updateData := &model.Student{
		ID: id,
	}
	if req.ProgramStudy != "" {
		updateData.ProgramStudy = req.ProgramStudy
	}
	if req.AcademicYear != "" {
		updateData.AcademicYear = req.AcademicYear
	}
	if req.AdvisorID != "" {
		updateData.AdvisorID = req.AdvisorID
	}

	if err := s.studentRepo.UpdateStudent(updateData); err != nil {
		return helper.NewApiError(500, "gagal mengupdate student: "+err.Error())
	}

	return nil
}

// DeleteStudent menghapus student
func (s *studentServiceImpl) DeleteStudent(id string) *helper.ApiError {
	if id == "" {
		return helper.NewApiError(400, "id tidak boleh kosong")
	}

	// Cek student ada atau tidak
	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return helper.NewApiError(500, "gagal mengambil student: "+err.Error())
	}
	if student == nil {
		return helper.NewApiError(404, "student tidak ditemukan")
	}

	if err := s.studentRepo.DeleteStudent(id); err != nil {
		return helper.NewApiError(500, "gagal menghapus student: "+err.Error())
	}

	return nil
}

// HandleGetAll menangani GET /api/students
func (s *studentServiceImpl) HandleGetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	data, apiErr := s.GetAllStudents(page, pageSize)
	if apiErr != nil {
		return c.Status(apiErr.Code).JSON(fiber.Map{
			"success": false,
			"message": apiErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// HandleGetByID menangani GET /api/students/:id
func (s *studentServiceImpl) HandleGetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	data, apiErr := s.GetStudentByID(id)
	if apiErr != nil {
		return c.Status(apiErr.Code).JSON(fiber.Map{
			"success": false,
			"message": apiErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// HandleCreate menangani POST /api/students
func (s *studentServiceImpl) HandleCreate(c *fiber.Ctx) error {
	var req model.StudentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	data, apiErr := s.CreateStudent(req)
	if apiErr != nil {
		return c.Status(apiErr.Code).JSON(fiber.Map{
			"success": false,
			"message": apiErr.Message,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// HandleUpdate menangani PUT /api/students/:id
func (s *studentServiceImpl) HandleUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.StudentUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	apiErr := s.UpdateStudent(id, req)
	if apiErr != nil {
		return c.Status(apiErr.Code).JSON(fiber.Map{
			"success": false,
			"message": apiErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "student updated successfully",
	})
}

// HandleDelete menangani DELETE /api/students/:id
func (s *studentServiceImpl) HandleDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	apiErr := s.DeleteStudent(id)
	if apiErr != nil {
		return c.Status(apiErr.Code).JSON(fiber.Map{
			"success": false,
			"message": apiErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "student deleted successfully",
	})
}
