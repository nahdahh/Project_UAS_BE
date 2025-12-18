package service

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetAllAchievements_Success tests getting all achievements successfully
func TestGetAllAchievements_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
		Points:          0, // Points will be assigned during verification
	}
	mockAchRepo.Create(achievement, studentID)
	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return service.GetAllAchievements(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements?page=1&page_size=10", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
	var response model.APIResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	assert.Equal(t, "success", response.Status)
	assert.Contains(t, response.Message, "achievements berhasil diambil")
}

// TestGetAllAchievements_StudentRole tests student can only see their achievements
func TestGetAllAchievements_StudentRole(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	studentID := uuid.New().String()
	userID := uuid.New().String()
	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
	})
	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.GetAllAchievements(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestGetAllAchievements_StudentNotFound tests student not found scenario
func TestGetAllAchievements_StudentNotFound(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	userID := uuid.New().String()
	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.GetAllAchievements(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 404, resp.StatusCode)
}

// TestGetAchievementDetail_Success tests getting achievement detail successfully
func TestGetAchievementDetail_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
		Points:          0, // Points will be assigned during verification
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	app.Get("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.GetAchievementDetail(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements/"+achWithRef.StudentID, nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestGetAchievementDetail_NotFound tests achievement not found scenario
func TestGetAchievementDetail_NotFound(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	userID := uuid.New().String()
	app.Get("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return service.GetAchievementDetail(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements/nonexistent", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 404, resp.StatusCode)
}

// TestGetAchievementDetail_Unauthorized tests unauthorized access
func TestGetAchievementDetail_Unauthorized(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	studentID := uuid.New().String()
	userID := uuid.New().String()
	otherUserID := uuid.New().String()

	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
	})
	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Title:           "Test Achievement",
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	app.Get("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", otherUserID)
		c.Locals("role", "Mahasiswa")
		return service.GetAchievementDetail(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements/"+achWithRef.StudentID, nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 401, resp.StatusCode)
}

// TestCreateAchievement_Success tests creating achievement successfully
func TestCreateAchievement_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	studentID := uuid.New().String()
	userID := uuid.New().String()
	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
	})
	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.CreateAchievement(c)
	})
	reqBody := model.CreateAchievementRequest{
		AchievementType: "academic",
		Title:           "New Achievement",
		Description:     "Test Description",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 201, resp.StatusCode)
	var response model.APIResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	assert.Equal(t, "success", response.Status)
	assert.Contains(t, response.Message, "achievement berhasil dibuat")
}

// TestCreateAchievement_InvalidJSON tests creating achievement with invalid JSON
func TestCreateAchievement_InvalidJSON(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	userID := uuid.New().String()
	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.CreateAchievement(c)
	})
	// Act
	req := httptest.NewRequest("POST", "/achievements", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestCreateAchievement_StudentNotFound tests student not found scenario
func TestCreateAchievement_StudentNotFound(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	userID := uuid.New().String()
	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.CreateAchievement(c)
	})
	reqBody := model.CreateAchievementRequest{
		AchievementType: "academic",
		Title:           "New Achievement",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 404, resp.StatusCode)
}

// TestUpdateAchievement_Success tests updating achievement successfully
func TestUpdateAchievement_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
		Title:           "Old Title",
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	app.Put("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.UpdateAchievement(c)
	})
	newTitle := "Updated Title"
	reqBody := model.UpdateAchievementRequest{
		Title: &newTitle,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("PUT", "/achievements/"+achievementID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestUpdateAchievement_NotDraft tests updating non-draft achievement
func TestUpdateAchievement_NotDraft(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID

	// Submit the achievement to change status
	mockAchRepo.Submit(achievementID)
	app.Put("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.UpdateAchievement(c)
	})
	newTitle := "Updated Title"
	reqBody := model.UpdateAchievementRequest{
		Title: &newTitle,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("PUT", "/achievements/"+achievementID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestSubmitAchievement_Success tests submitting achievement successfully
func TestSubmitAchievement_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.SubmitAchievement(c)
	})
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/submit", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestSubmitAchievement_NotDraft tests submitting non-draft achievement
func TestSubmitAchievement_NotDraft(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	mockAchRepo.Submit(achievementID)
	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.SubmitAchievement(c)
	})
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/submit", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestVerifyAchievement_Success tests verifying achievement successfully
func TestVerifyAchievement_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	studentID := uuid.New().String()
	userID := uuid.New().String()
	lecturerID := uuid.New().String()
	lecturerUserID := uuid.New().String()
	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
		AdvisorID: lecturerID,
	})
	mockLecturerRepo.CreateLecturer(&model.Lecturer{
		ID:         lecturerID,
		UserID:     lecturerUserID,
		LecturerID: "789012",
		Department: "Computer Science",
	})
	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Title:           "Test Achievement",
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	mockAchRepo.Submit(achievementID)
	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("userID", lecturerUserID)
		c.Locals("role", "Dosen Wali")
		return service.VerifyAchievement(c)
	})
	reqBody := model.VerifyAchievementRequest{
		Points: 100,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestVerifyAchievement_StudentCannotVerify tests student cannot verify
func TestVerifyAchievement_StudentCannotVerify(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	mockAchRepo.Submit(achievementID)
	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.VerifyAchievement(c)
	})
	reqBody := model.VerifyAchievementRequest{
		Points: 100,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 403, resp.StatusCode)
}

// TestVerifyAchievement_NegativePoints tests negative points validation
func TestVerifyAchievement_NegativePoints(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	userID := uuid.New().String()
	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return service.VerifyAchievement(c)
	})
	reqBody := model.VerifyAchievementRequest{
		Points: -10,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements/test-id/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestRejectAchievement_Success tests rejecting achievement successfully
func TestRejectAchievement_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	studentID := uuid.New().String()
	userID := uuid.New().String()
	lecturerID := uuid.New().String()
	lecturerUserID := uuid.New().String()
	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
		AdvisorID: lecturerID,
	})
	mockLecturerRepo.CreateLecturer(&model.Lecturer{
		ID:         lecturerID,
		UserID:     lecturerUserID,
		LecturerID: "789012",
		Department: "Computer Science",
	})
	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Title:           "Test Achievement",
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	mockAchRepo.Submit(achievementID)
	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		c.Locals("userID", lecturerUserID)
		c.Locals("role", "Dosen Wali")
		return service.RejectAchievement(c)
	})
	reqBody := map[string]string{
		"rejection_note": "Needs more documentation",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/reject", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestRejectAchievement_EmptyNote tests rejecting without note
func TestRejectAchievement_EmptyNote(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	mockAchRepo.Submit(achievementID)
	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Admin")
		return service.RejectAchievement(c)
	})
	reqBody := map[string]string{
		"rejection_note": "",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/reject", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestDeleteAchievement_Success tests deleting achievement successfully
func TestDeleteAchievement_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	app.Delete("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.DeleteAchievement(c)
	})
	// Act
	req := httptest.NewRequest("DELETE", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestDeleteAchievement_NotDraft tests deleting non-draft achievement
func TestDeleteAchievement_NotDraft(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	mockAchRepo.Submit(achievementID)
	app.Delete("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.DeleteAchievement(c)
	})
	// Act
	req := httptest.NewRequest("DELETE", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestGetAchievementHistory_Success tests getting achievement history successfully
func TestGetAchievementHistory_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	app.Get("/achievements/:id/history", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.GetAchievementHistory(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements/"+achievementID+"/history", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestUploadAttachment_Success tests uploading attachment successfully
func TestUploadAttachment_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	app.Post("/achievements/:id/attachments", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.UploadAttachment(c)
	})
	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("test file content"))
	writer.Close()
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/attachments", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req, -1) // -1 timeout for file operations
	// Assert
	assert.Equal(t, 201, resp.StatusCode)
}

// TestUploadAttachment_NoFile tests uploading without file
func TestUploadAttachment_NoFile(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
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
	}
	achWithRef, _ := mockAchRepo.Create(achievement, studentID)
	achievementID := achWithRef.StudentID
	app.Post("/achievements/:id/attachments", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.UploadAttachment(c)
	})
	// Act
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/attachments", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 400, resp.StatusCode)
}

// TestGetAdviseeAchievements_Success tests getting advisee achievements
func TestGetAdviseeAchievements_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	lecturerID := uuid.New().String()
	studentID := uuid.New().String()
	userID := uuid.New().String()
	mockStudentRepo.CreateStudent(&model.Student{
		ID:        studentID,
		UserID:    userID,
		StudentID: "123456",
		AdvisorID: lecturerID,
	})
	app.Get("/achievements/advisee/list", func(c *fiber.Ctx) error {
		c.Locals("userID", lecturerID)
		c.Locals("role", "Dosen Wali")
		return service.GetAdviseeAchievements(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements/advisee/list", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 200, resp.StatusCode)
}

// TestGetAdviseeAchievements_StudentForbidden tests student cannot access
func TestGetAdviseeAchievements_StudentForbidden(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAchRepo := repository.NewMockAchievementRepository()
	mockStudentRepo := repository.NewMockStudentRepository()
	mockLecturerRepo := repository.NewMockLecturerRepository()
	service := NewAchievementService(mockAchRepo, mockStudentRepo, mockLecturerRepo)
	userID := uuid.New().String()
	app.Get("/achievements/advisee/list", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		c.Locals("role", "Mahasiswa")
		return service.GetAdviseeAchievements(c)
	})
	// Act
	req := httptest.NewRequest("GET", "/achievements/advisee/list", nil)
	resp, _ := app.Test(req)
	// Assert
	assert.Equal(t, 403, resp.StatusCode)
}
