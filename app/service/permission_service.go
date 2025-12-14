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
