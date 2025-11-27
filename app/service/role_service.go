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
	// CreateRole membuat role baru
	CreateRole(name, description string) (*model.Role, error)

	// GetRoleByID mengambil role berdasarkan ID
	GetRoleByID(id string) (*model.Role, error)

	// GetRoleByName mengambil role berdasarkan nama
	GetRoleByName(name string) (*model.Role, error)

	// GetAllRoles mengambil semua role
	GetAllRoles() ([]*model.Role, error)

	// UpdateRole mengubah data role
	UpdateRole(id, description string) error

	// AssignPermissionToRole menambahkan permission ke role
	AssignPermissionToRole(roleID, permissionID string) error

	// RemovePermissionFromRole menghapus permission dari role
	RemovePermissionFromRole(roleID, permissionID string) error

	// HTTP Handler methods
	HandleGetAllRoles(c *fiber.Ctx) error
	HandleGetRoleByID(c *fiber.Ctx) error
	HandleCreateRole(c *fiber.Ctx) error
	HandleUpdateRole(c *fiber.Ctx) error
	HandleAssignPermissionToRole(c *fiber.Ctx) error
	HandleRemovePermissionFromRole(c *fiber.Ctx) error
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

// CreateRole membuat role baru
func (s *roleServiceImpl) CreateRole(name, description string) (*model.Role, error) {
	// Validasi input
	if name == "" {
		return nil, errors.New("nama role tidak boleh kosong")
	}

	// Cek apakah role sudah ada
	existing, err := s.roleRepo.GetRoleByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("nama role sudah terdaftar")
	}

	// Buat role baru
	role := &model.Role{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}

	// Simpan ke database
	err = s.roleRepo.CreateRole(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID mengambil role berdasarkan ID
func (s *roleServiceImpl) GetRoleByID(id string) (*model.Role, error) {
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}

	return role, nil
}

// GetRoleByName mengambil role berdasarkan nama
func (s *roleServiceImpl) GetRoleByName(name string) (*model.Role, error) {
	role, err := s.roleRepo.GetRoleByName(name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}

	return role, nil
}

// GetAllRoles mengambil semua role
func (s *roleServiceImpl) GetAllRoles() ([]*model.Role, error) {
	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// UpdateRole mengubah data role
func (s *roleServiceImpl) UpdateRole(id, description string) error {
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role tidak ditemukan")
	}

	// Update field
	if description != "" {
		role.Description = description
	}

	return s.roleRepo.UpdateRole(role)
}

// AssignPermissionToRole menambahkan permission ke role
func (s *roleServiceImpl) AssignPermissionToRole(roleID, permissionID string) error {
	// Validasi role dan permission ada
	role, err := s.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role tidak ditemukan")
	}

	permission, err := s.permissionRepo.GetPermissionByID(permissionID)
	if err != nil {
		return err
	}
	if permission == nil {
		return errors.New("permission tidak ditemukan")
	}

	// Assign permission ke role
	return s.roleRepo.AssignPermissionToRole(roleID, permissionID)
}

// RemovePermissionFromRole menghapus permission dari role
func (s *roleServiceImpl) RemovePermissionFromRole(roleID, permissionID string) error {
	// Validasi role ada
	role, err := s.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role tidak ditemukan")
	}

	// Remove permission dari role
	return s.roleRepo.RemovePermissionFromRole(roleID, permissionID)
}

// ===== HTTP HANDLER METHODS (Route Layer) =====
// Handlers hanya parse input → panggil business logic → return response

// HandleGetAllRoles menangani GET /api/v1/roles
func (s *roleServiceImpl) HandleGetAllRoles(c *fiber.Ctx) error {
	roles, err := s.GetAllRoles()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "gagal mengambil data role: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data role berhasil diambil", roles)
}

// HandleGetRoleByID menangani GET /api/v1/roles/:id
func (s *roleServiceImpl) HandleGetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	role, err := s.GetRoleByID(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "role tidak ditemukan: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "data role berhasil diambil", role)
}

// HandleCreateRole menangani POST /api/v1/roles
func (s *roleServiceImpl) HandleCreateRole(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
	}

	role, err := s.CreateRole(req.Name, req.Description)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal membuat role: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "role berhasil dibuat", role)
}

// HandleUpdateRole menangani PUT /api/v1/roles/:id
func (s *roleServiceImpl) HandleUpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "id tidak boleh kosong")
	}

	var req struct {
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "format request tidak valid")
	}

	err := s.UpdateRole(id, req.Description)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal mengupdate role: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "role berhasil diupdate", nil)
}

// HandleAssignPermissionToRole menangani POST /api/v1/roles/:id/permissions/:permissionId
func (s *roleServiceImpl) HandleAssignPermissionToRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	permissionID := c.Params("permissionId")

	if roleID == "" || permissionID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "role id dan permission id tidak boleh kosong")
	}

	err := s.AssignPermissionToRole(roleID, permissionID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal assign permission: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "permission berhasil di-assign ke role", nil)
}

// HandleRemovePermissionFromRole menangani DELETE /api/v1/roles/:id/permissions/:permissionId
func (s *roleServiceImpl) HandleRemovePermissionFromRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	permissionID := c.Params("permissionId")

	if roleID == "" || permissionID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "role id dan permission id tidak boleh kosong")
	}

	err := s.RemovePermissionFromRole(roleID, permissionID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "gagal remove permission: "+err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "permission berhasil dihapus dari role", nil)
}
