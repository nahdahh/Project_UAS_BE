package service

import (
	"strconv"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	AssignRole(c *fiber.Ctx) error
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// CreateUser godoc
// @Summary Buat user baru
// @Description Membuat user baru dengan role tertentu
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body model.CreateUserRequest true "Data user"
// @Success 201 {object} model.APIResponse{data=model.User} "User berhasil dibuat"
// @Failure 400 {object} model.APIResponse "Format request tidak valid atau username/email sudah terdaftar"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /users [post]
func (s *userServiceImpl) CreateUser(c *fiber.Ctx) error {
	req := new(model.CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	// Validasi input
	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "semua field harus diisi",
		})
	}

	// Cek username unik
	existingUser, _ := s.userRepo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "username sudah terdaftar",
		})
	}

	// Cek email unik
	existingUser, _ = s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "email sudah terdaftar",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal hash password",
		})
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal membuat user: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "user berhasil dibuat",
		Data:    user,
	})
}

// GetUserByID godoc
// @Summary Dapatkan user berdasarkan ID
// @Description Mengambil data user berdasarkan ID beserta permissions
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} model.APIResponse{data=model.UserWithRole} "Data user berhasil diambil"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "User tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /users/{id} [get]
func (s *userServiceImpl) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil user: " + err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "user tidak ditemukan",
		})
	}

	permissions, err := s.userRepo.GetUserPermissions(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil permissions: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "user berhasil diambil",
		Data: &model.UserWithRole{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			IsActive:    user.IsActive,
			Permissions: permissions,
			CreatedAt:   user.CreatedAt.String(),
		},
	})
}

// GetAllUsers godoc
// @Summary Dapatkan semua user
// @Description Mengambil daftar semua user dengan pagination
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman" default(1)
// @Param page_size query int false "Jumlah data per halaman" default(10)
// @Success 200 {object} model.APIResponse{data=object{users=[]model.User,pagination=object}} "Daftar user berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /users [get]
func (s *userServiceImpl) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.GetAllUsers(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil users: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "users berhasil diambil",
		Data: map[string]interface{}{
			"users": users,
			"pagination": map[string]interface{}{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + pageSize - 1) / pageSize,
			},
		},
	})
}

// UpdateUser godoc
// @Summary Update data user
// @Description Memperbarui data user (username, email, full_name)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body object{username=string,email=string,full_name=string} true "Data yang diupdate"
// @Success 200 {object} model.APIResponse{data=model.User} "User berhasil diupdate"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "User tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /users/{id} [put]
func (s *userServiceImpl) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	}

	req := new(UpdateUserRequest)
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

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil user: " + err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "user tidak ditemukan",
		})
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengupdate user: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "user berhasil diupdate",
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary Hapus user
// @Description Menghapus user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} model.APIResponse "User berhasil dihapus"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "User tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /users/{id} [delete]
func (s *userServiceImpl) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "id tidak boleh kosong",
		})
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil user: " + err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "user tidak ditemukan",
		})
	}

	if err := s.userRepo.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal menghapus user: " + err.Error(),
		})
	}

	// HTTP response
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "user berhasil dihapus",
	})
}

// AssignRole godoc
// @Summary Assign role ke user
// @Description Mengubah role user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body object{role_id=string} true "Role ID baru"
// @Success 200 {object} model.APIResponse "Role berhasil diassign"
// @Failure 400 {object} model.APIResponse "User ID dan role ID harus diisi"
// @Failure 404 {object} model.APIResponse "User tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /users/{id}/role [put]
func (s *userServiceImpl) AssignRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	type AssignRoleRequest struct {
		RoleID string `json:"role_id"`
	}

	req := new(AssignRoleRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "format request tidak valid: " + err.Error(),
		})
	}

	if userID == "" || req.RoleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Status:  "error",
			Message: "user_id dan role_id harus diisi",
		})
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal mengambil user: " + err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.APIResponse{
			Status:  "error",
			Message: "user tidak ditemukan",
		})
	}

	if err := s.userRepo.AssignRole(userID, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Status:  "error",
			Message: "gagal assign role: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "role berhasil diassign",
	})
}
