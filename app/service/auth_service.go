package service

import (
	"time"
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"
	"uas_be/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authServiceImpl struct {
	userRepo       repository.UserRepository
	permissionRepo repository.PermissionRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	permissionRepo repository.PermissionRepository,
) AuthService {
	return &authServiceImpl{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *authServiceImpl) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	if req.Username == "" || req.Password == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "username dan password harus diisi")
	}

	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil user")
	}
	if user == nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "username atau password salah")
	}

	if !user.IsActive {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "user tidak aktif")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "username atau password salah")
	}

	permissionNames, err := s.permissionRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil permissions")
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID, permissionNames, 24*time.Hour)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal generate token")
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.RoleID, permissionNames, 7*24*time.Hour)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal generate refresh token")
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "login berhasil",
		Data: map[string]interface{}{
			"token":         token,
			"refresh_token": refreshToken,
			"user": map[string]interface{}{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"role_id":   user.RoleID,
			},
			"permissions": permissionNames,
		},
	})
}

func (s *authServiceImpl) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "user tidak terautentikasi")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil user")
	}
	if user == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "user tidak ditemukan")
	}

	permissionNames, err := s.permissionRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil permissions")
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "profile berhasil diambil",
		Data: map[string]interface{}{
			"id":          user.ID,
			"username":    user.Username,
			"email":       user.Email,
			"full_name":   user.FullName,
			"role_id":     user.RoleID,
			"permissions": permissionNames,
			"is_active":   user.IsActive,
			"created_at":  user.CreatedAt.Format(time.RFC3339),
			"updated_at":  user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func (s *authServiceImpl) RefreshToken(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	req := new(RefreshRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	if req.RefreshToken == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "refresh_token harus diisi")
	}

	claims, err := utils.GetClaimsFromToken(req.RefreshToken)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "refresh token tidak valid: "+err.Error())
	}

	user, err := s.userRepo.GetUserByID(claims.Sub)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil user")
	}
	if user == nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "user tidak ditemukan")
	}

	if !user.IsActive {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "user tidak aktif")
	}

	permissionNames, err := s.permissionRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil permissions")
	}

	newToken, err := utils.GenerateToken(user.ID, user.RoleID, permissionNames, 24*time.Hour)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal generate token")
	}

	newRefreshToken, err := utils.GenerateToken(user.ID, user.RoleID, permissionNames, 7*24*time.Hour)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal generate refresh token")
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "token berhasil di-refresh",
		Data: map[string]interface{}{
			"token":         newToken,
			"refresh_token": newRefreshToken,
		},
	})
}

func (s *authServiceImpl) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "logout berhasil",
		Data:    nil,
	})
}
