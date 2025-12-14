package service

import (
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	GetStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error
}

type reportServiceImpl struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
}

func NewReportService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
) ReportService {
	return &reportServiceImpl{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
	}
}

func (s *reportServiceImpl) GetStatistics(c *fiber.Ctx) error {
	// Get all achievements
	achievements, _, err := s.achievementRepo.GetAllAchievements(1, 10000)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil data achievements",
		})
	}

	// Calculate statistics
	stats := make(map[string]interface{})
	totalAchievements := len(achievements)

	// Count by status
	statusCount := make(map[string]int)
	typeCount := make(map[string]int)
	totalPoints := 0

	for _, ach := range achievements {
		statusCount[ach.Status]++
		typeCount[ach.AchievementType]++
		if ach.Status == "verified" {
			totalPoints += ach.Points
		}
	}

	// Get top students (simplified - in production use more sophisticated query)
	studentAchievements := make(map[string]int)
	studentPoints := make(map[string]int)

	for _, ach := range achievements {
		if ach.Status == "verified" {
			studentAchievements[ach.StudentID]++
			studentPoints[ach.StudentID] += ach.Points
		}
	}

	stats["total_achievements"] = totalAchievements
	stats["by_status"] = statusCount
	stats["by_type"] = typeCount
	stats["total_points"] = totalPoints
	stats["verified_achievements"] = statusCount["verified"]
	stats["pending_achievements"] = statusCount["submitted"]

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "statistik berhasil diambil",
		Data:    stats,
	})
}

func (s *reportServiceImpl) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")

	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "student_id tidak boleh kosong",
		})
	}

	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil student",
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
			Message: "gagal mengambil achievements",
		})
	}

	// Calculate student statistics
	report := make(map[string]interface{})
	totalAchievements := len(achievements)
	totalPoints := 0
	statusCount := make(map[string]int)
	typeCount := make(map[string]int)

	for _, ach := range achievements {
		statusCount[ach.Status]++
		typeCount[ach.AchievementType]++
		if ach.Status == "verified" {
			totalPoints += ach.Points
		}
	}

	report["student"] = student
	report["total_achievements"] = totalAchievements
	report["total_points"] = totalPoints
	report["by_status"] = statusCount
	report["by_type"] = typeCount
	report["achievements"] = achievements

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "report berhasil diambil",
		Data:    report,
	})
}
