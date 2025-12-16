package service

import (
	"strconv"
	"time"
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"

	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	GetStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error
	GetStatisticsByPeriod(c *fiber.Ctx) error
	GetStatisticsByType(c *fiber.Ctx) error
	GetTopStudents(c *fiber.Ctx) error
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

// GetStatistics godoc
// @Summary Dapatkan statistik prestasi
// @Description Mengambil statistik prestasi berdasarkan role user
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=object} "Statistik berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /reports/statistics [get]
func (s *reportServiceImpl) GetStatistics(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	var achievements []*model.AchievementWithReference
	var err error

	// Apply RBAC: Admin sees all, Dosen Wali sees advisees, Mahasiswa sees own
	if role == "Admin" {
		achievements, _, err = s.achievementRepo.GetAllAchievements(1, 10000)
	} else if role == "Dosen Wali" {
		students, err := s.studentRepo.GetStudentsByAdvisorID(userID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data mahasiswa")
		}

		for _, student := range students {
			studentAchs, err := s.achievementRepo.GetAchievementsByStudentID(student.ID)
			if err == nil {
				achievements = append(achievements, studentAchs...)
			}
		}
	} else {
		achievements, err = s.achievementRepo.GetAchievementsByStudentID(userID)
	}

	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data achievements")
	}

	stats := make(map[string]interface{})
	totalAchievements := len(achievements)

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

	stats["total_achievements"] = totalAchievements
	stats["by_status"] = statusCount
	stats["by_type"] = typeCount
	stats["total_points"] = totalPoints
	stats["verified_achievements"] = statusCount["verified"]
	stats["pending_achievements"] = statusCount["submitted"]
	stats["draft_achievements"] = statusCount["draft"]
	stats["rejected_achievements"] = statusCount["rejected"]

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "statistik berhasil diambil",
		Data:    stats,
	})
}

// GetStudentReport godoc
// @Summary Dapatkan laporan mahasiswa
// @Description Mengambil laporan lengkap prestasi mahasiswa berdasarkan ID
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} model.APIResponse{data=object} "Laporan berhasil diambil"
// @Failure 400 {object} model.APIResponse "Student ID tidak boleh kosong"
// @Failure 403 {object} model.APIResponse "Anda tidak memiliki akses ke report ini"
// @Failure 404 {object} model.APIResponse "Mahasiswa tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /reports/student/{id} [get]
func (s *reportServiceImpl) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	if studentID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "student_id tidak boleh kosong")
	}

	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil student")
	}
	if student == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	// RBAC check
	if role == "Mahasiswa" && studentID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "anda tidak memiliki akses ke report ini")
	}

	if role == "Dosen Wali" && student.AdvisorID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "student ini bukan bimbingan anda")
	}

	achievements, err := s.achievementRepo.GetAchievementsByStudentID(studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil achievements")
	}

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

// GetStatisticsByPeriod godoc
// @Summary Dapatkan statistik berdasarkan periode
// @Description Mengambil statistik prestasi dalam rentang waktu tertentu
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Tanggal mulai (YYYY-MM-DD)"
// @Param end_date query string true "Tanggal akhir (YYYY-MM-DD)"
// @Success 200 {object} model.APIResponse{data=object} "Statistik berdasarkan periode berhasil diambil"
// @Failure 400 {object} model.APIResponse "Start date dan end date harus diisi atau format tidak valid"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /reports/statistics/period [get]
func (s *reportServiceImpl) GetStatisticsByPeriod(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "start_date dan end_date harus diisi")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format start_date tidak valid (gunakan YYYY-MM-DD)")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format end_date tidak valid (gunakan YYYY-MM-DD)")
	}

	stats, err := s.achievementRepo.GetAchievementStatsByPeriod(startDate, endDate, role, userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil statistik: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "statistik berdasarkan periode berhasil diambil",
		Data:    stats,
	})
}

// GetStatisticsByType godoc
// @Summary Dapatkan statistik berdasarkan tipe
// @Description Mengambil statistik prestasi berdasarkan tipe prestasi
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=object} "Statistik berdasarkan tipe berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /reports/statistics/type [get]
func (s *reportServiceImpl) GetStatisticsByType(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	stats, err := s.achievementRepo.GetAchievementStatsByType(role, userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil statistik: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "statistik berdasarkan tipe berhasil diambil",
		Data:    stats,
	})
}

// GetTopStudents godoc
// @Summary Dapatkan mahasiswa dengan prestasi terbaik
// @Description Mengambil daftar mahasiswa dengan poin prestasi tertinggi (hanya admin)
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Jumlah mahasiswa yang ditampilkan" default(10)
// @Success 200 {object} model.APIResponse{data=object{top_students=[]object,limit=int}} "Top students berhasil diambil"
// @Failure 403 {object} model.APIResponse "Hanya admin yang dapat mengakses data ini"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /reports/top-students [get]
func (s *reportServiceImpl) GetTopStudents(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "Admin" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "hanya admin yang dapat mengakses data ini")
	}

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	topStudents, err := s.achievementRepo.GetTopStudents(limit)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil top students: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "top students berhasil diambil",
		Data: map[string]interface{}{
			"top_students": topStudents,
			"limit":        limit,
		},
	})
}
