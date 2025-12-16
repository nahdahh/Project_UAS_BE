package service

import (
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PermissionService interface {
	GetAllPermissions(c *fiber.Ctx) error
	CreatePermission(c *fiber.Ctx) error
	GetPermissionsByRoleID(c *fiber.Ctx) error
}

type permissionServiceImpl struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionServiceImpl{
		permissionRepo: permissionRepo,
	}
}

// GetAllPermissions godoc
// @Summary Dapatkan semua permission
// @Description Mengambil daftar semua permission yang tersedia
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=[]model.Permission} "Daftar permission berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /permissions [get]
func (s *permissionServiceImpl) GetAllPermissions(c *fiber.Ctx) error {
	permissions, err := s.permissionRepo.GetAllPermissions()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil permissions: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "permissions berhasil diambil",
		Data:    permissions,
	})
}

// CreatePermission godoc
// @Summary Buat permission baru
// @Description Membuat permission baru untuk sistem RBAC
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{name=string,resource=string,action=string,description=string} true "Data permission"
// @Success 201 {object} model.APIResponse{data=model.Permission} "Permission berhasil dibuat"
// @Failure 400 {object} model.APIResponse "Format request tidak valid atau field wajib tidak diisi"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /permissions [post]
func (s *permissionServiceImpl) CreatePermission(c *fiber.Ctx) error {
	type CreatePermissionRequest struct {
		Name        string `json:"name"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
		Description string `json:"description"`
	}

	req := new(CreatePermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	if req.Name == "" || req.Resource == "" || req.Action == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "name, resource, dan action harus diisi")
	}

	permission := &model.Permission{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := s.permissionRepo.CreatePermission(permission); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat permission: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Status:  "success",
		Message: "permission berhasil dibuat",
		Data:    permission,
	})
}

// GetPermissionsByRoleID godoc
// @Summary Dapatkan permission berdasarkan role
// @Description Mengambil daftar permission yang dimiliki oleh role tertentu
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roleId path string true "Role ID"
// @Success 200 {object} model.APIResponse{data=[]string} "Daftar permission berhasil diambil"
// @Failure 400 {object} model.APIResponse "Role ID tidak boleh kosong"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /permissions/role/{roleId} [get]
func (s *permissionServiceImpl) GetPermissionsByRoleID(c *fiber.Ctx) error {
	roleID := c.Params("roleId")

	if roleID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "roleId tidak boleh kosong")
	}

	permissions, err := s.permissionRepo.GetPermissionsByRoleID(roleID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil permissions: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Status:  "success",
		Message: "permissions berhasil diambil",
		Data:    permissions,
	})
}
