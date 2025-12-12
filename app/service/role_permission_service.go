package service

import (
	"uas_be/app/model"
	"uas_be/app/repository"
)

type RolePermissionService struct {
	repo *repository.RolePermissionRepository
}

func NewRolePermissionService(repo *repository.RolePermissionRepository) *RolePermissionService {
	return &RolePermissionService{repo: repo}
}

func (s *RolePermissionService) AssignPermission(roleID, permissionID string) error {
	return s.repo.AssignPermission(roleID, permissionID)
}

func (s *RolePermissionService) RemovePermission(roleID, permissionID string) error {
	return s.repo.RemovePermission(roleID, permissionID)
}

func (s *RolePermissionService) GetPermissions(roleID string) ([]model.Permission, error) {
	return s.repo.GetPermissionsByRole(roleID)
}
