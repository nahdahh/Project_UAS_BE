package service

import (
	"net/http/httptest"
	"testing"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReportService_GetStatistics(t *testing.T) {
	// Setup mocks
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()

	service := NewReportService(mockAchRepo, mockStudentRepo, mockLecturerRepo)

	// Setup test data
	studentID := uuid.New().String()
	userID := uuid.New().String()

	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
	})

	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          10,
	}
	mockAchRepo.Create(achievement, studentID)

	// Test
	app := fiber.New()

	// Add middleware to set locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/statistics", service.GetStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReportService_GetStudentReport(t *testing.T) {
	// Setup mocks
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()

	service := NewReportService(mockAchRepo, mockStudentRepo, mockLecturerRepo)

	// Setup test data
	studentID := uuid.New().String()
	userID := uuid.New().String()

	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
	})

	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          10,
	}
	mockAchRepo.Create(achievement, studentID)

	// Test
	app := fiber.New()

	// Add middleware to set locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/student-report/:id", service.GetStudentReport)

	req := httptest.NewRequest("GET", "/student-report/"+studentID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReportService_GetStatisticsByPeriod(t *testing.T) {
	// Setup mocks
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()

	service := NewReportService(mockAchRepo, mockStudentRepo, mockLecturerRepo)

	// Test
	app := fiber.New()

	// Add middleware to set locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/statistics/period", service.GetStatisticsByPeriod)

	req := httptest.NewRequest("GET", "/statistics/period?start_date=2024-01-01&end_date=2024-12-31", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReportService_GetStatisticsByType(t *testing.T) {
	// Setup mocks
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()

	service := NewReportService(mockAchRepo, mockStudentRepo, mockLecturerRepo)

	// Test
	app := fiber.New()

	// Add middleware to set locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/statistics/type", service.GetStatisticsByType)

	req := httptest.NewRequest("GET", "/statistics/type?achievement_type=academic", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReportService_GetTopStudents(t *testing.T) {
	// Setup mocks
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()

	service := NewReportService(mockAchRepo, mockStudentRepo, mockLecturerRepo)

	// Setup test data
	studentID := uuid.New().String()
	userID := uuid.New().String()

	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
	})

	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          10,
	}
	mockAchRepo.Create(achievement, studentID)

	// Test
	app := fiber.New()

	// Add middleware to set locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/top-students", service.GetTopStudents)

	req := httptest.NewRequest("GET", "/top-students?limit=10", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}