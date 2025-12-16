package service

import (
	"errors"
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoleService adalah interface untuk business logic role/peran
type RoleService interface {
	CreateRole(c *fiber.Ctx) error
	GetRoleByID(c *fiber.Ctx) error
	GetRoleByName(c *fiber.Ctx) error
	GetAllRoles(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	AssignPermission(c *fiber.Ctx) error
	RemovePermission(c *fiber.Ctx) error
}

// roleServiceImpl adalah implementasi dari RoleService
type roleServiceImpl struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewRoleService membuat instance service role baru
func NewRoleService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) RoleService {
	return &roleServiceImpl{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// CreateRole godoc
// @Summary Buat role baru
// @Description Membuat role baru untuk sistem RBAC
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{name=string,description=string} true "Data role"
// @Success 201 {object} model.APIResponse{data=model.Role} "Role berhasil dibuat"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles [post]
func (s *roleServiceImpl) CreateRole(c *fiber.Ctx) error {
	type CreateRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	req := new(CreateRoleRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	role, err := s.createRole(req.Name, req.Description)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "role berhasil dibuat", role)
}

// GetRoleByID godoc
// @Summary Dapatkan role berdasarkan ID
// @Description Mengambil data role berdasarkan ID
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} model.APIResponse{data=model.Role} "Data role berhasil diambil"
// @Failure 400 {object} model.APIResponse "ID tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Role tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles/{id} [get]
func (s *roleServiceImpl) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	role, err := s.getRoleByID(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	if role == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "role tidak ditemukan")
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "role berhasil diambil", role)
}

// GetRoleByName godoc
// @Summary Dapatkan role berdasarkan nama
// @Description Mengambil data role berdasarkan nama
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string true "Role Name"
// @Success 200 {object} model.APIResponse{data=model.Role} "Data role berhasil diambil"
// @Failure 400 {object} model.APIResponse "Name tidak boleh kosong"
// @Failure 404 {object} model.APIResponse "Role tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles/name/{name} [get]
func (s *roleServiceImpl) GetRoleByName(c *fiber.Ctx) error {
	name := c.Params("name")

	role, err := s.getRoleByName(name)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	if role == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "role tidak ditemukan")
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "role berhasil diambil", role)
}

// GetAllRoles godoc
// @Summary Dapatkan semua role
// @Description Mengambil daftar semua role yang tersedia
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=[]model.Role} "Daftar role berhasil diambil"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles [get]
func (s *roleServiceImpl) GetAllRoles(c *fiber.Ctx) error {
	roles, err := s.getAllRoles()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "roles berhasil diambil", roles)
}

// UpdateRole godoc
// @Summary Update role
// @Description Memperbarui deskripsi role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Param body body object{description=string} true "Deskripsi baru"
// @Success 200 {object} model.APIResponse{data=model.Role} "Role berhasil diupdate"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "Role tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles/{id} [put]
func (s *roleServiceImpl) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateRoleRequest struct {
		Description string `json:"description"`
	}

	req := new(UpdateRoleRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	role, err := s.updateRole(id, req.Description)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "role berhasil diupdate", role)
}

// AssignPermission godoc
// @Summary Assign permission ke role
// @Description Menambahkan permission ke role tertentu
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{role_id=string,permission_id=string} true "Role ID dan Permission ID"
// @Success 200 {object} model.APIResponse "Permission berhasil ditambahkan ke role"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "Role atau permission tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles/assign-permission [post]
func (s *roleServiceImpl) AssignPermission(c *fiber.Ctx) error {
	type AssignPermissionRequest struct {
		RoleID       string `json:"role_id"`
		PermissionID string `json:"permission_id"`
	}

	req := new(AssignPermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	err := s.assignPermissionToRole(req.RoleID, req.PermissionID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "permission berhasil ditambahkan ke role", nil)
}

// RemovePermission godoc
// @Summary Hapus permission dari role
// @Description Menghapus permission dari role tertentu
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{role_id=string,permission_id=string} true "Role ID dan Permission ID"
// @Success 200 {object} model.APIResponse "Permission berhasil dihapus dari role"
// @Failure 400 {object} model.APIResponse "Format request tidak valid"
// @Failure 404 {object} model.APIResponse "Role tidak ditemukan"
// @Failure 500 {object} model.APIResponse "Internal server error"
// @Router /roles/remove-permission [post]
func (s *roleServiceImpl) RemovePermission(c *fiber.Ctx) error {
	type RemovePermissionRequest struct {
		RoleID       string `json:"role_id"`
		PermissionID string `json:"permission_id"`
	}

	req := new(RemovePermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid: "+err.Error())
	}

	err := s.removePermissionFromRole(req.RoleID, req.PermissionID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "permission berhasil dihapus dari role", nil)
}

func (s *roleServiceImpl) createRole(name, description string) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("nama role tidak boleh kosong")
	}

	existing, _ := s.roleRepo.GetRoleByName(name)
	if existing != nil {
		return nil, errors.New("nama role sudah terdaftar")
	}

	role := &model.Role{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}

	if err := s.roleRepo.CreateRole(role); err != nil {
		return nil, errors.New("gagal membuat role: " + err.Error())
	}

	return role, nil
}

func (s *roleServiceImpl) getRoleByID(id string) (*model.Role, error) {
	if id == "" {
		return nil, errors.New("id tidak boleh kosong")
	}

	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, errors.New("gagal mengambil role: " + err.Error())
	}

	return role, nil
}

func (s *roleServiceImpl) getRoleByName(name string) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("name tidak boleh kosong")
	}

	role, err := s.roleRepo.GetRoleByName(name)
	if err != nil {
		return nil, errors.New("gagal mengambil role: " + err.Error())
	}

	return role, nil
}

func (s *roleServiceImpl) getAllRoles() ([]*model.Role, error) {
	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return nil, errors.New("gagal mengambil roles: " + err.Error())
	}

	return roles, nil
}

func (s *roleServiceImpl) updateRole(id, description string) (*model.Role, error) {
	if id == "" {
		return nil, errors.New("id tidak boleh kosong")
	}

	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, errors.New("gagal mengambil role: " + err.Error())
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}

	if description != "" {
		role.Description = description
	}

	if err := s.roleRepo.UpdateRole(role); err != nil {
		return nil, errors.New("gagal mengupdate role: " + err.Error())
	}

	return role, nil
}

func (s *roleServiceImpl) assignPermissionToRole(roleID, permissionID string) error {
	if roleID == "" || permissionID == "" {
		return errors.New("role_id dan permission_id tidak boleh kosong")
	}

	role, _ := s.roleRepo.GetRoleByID(roleID)
	if role == nil {
		return errors.New("role tidak ditemukan")
	}

	permission, _ := s.permissionRepo.GetPermissionByID(permissionID)
	if permission == nil {
		return errors.New("permission tidak ditemukan")
	}

	if err := s.roleRepo.AssignPermissionToRole(roleID, permissionID); err != nil {
		return errors.New("gagal assign permission ke role: " + err.Error())
	}

	return nil
}

func (s *roleServiceImpl) removePermissionFromRole(roleID, permissionID string) error {
	if roleID == "" || permissionID == "" {
		return errors.New("role_id dan permission_id tidak boleh kosong")
	}

	role, _ := s.roleRepo.GetRoleByID(roleID)
	if role == nil {
		return errors.New("role tidak ditemukan")
	}

	if err := s.roleRepo.RemovePermissionFromRole(roleID, permissionID); err != nil {
		return errors.New("gagal remove permission dari role: " + err.Error())
	}

	return nil
}
