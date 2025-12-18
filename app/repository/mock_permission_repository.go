package repository

import (
	"errors"
	"uas_be/app/model"
)

type MockPermissionRepository struct {
	permissions map[string]*model.Permission
}

func NewMockPermissionRepository() *MockPermissionRepository {
	return &MockPermissionRepository{
		permissions: make(map[string]*model.Permission),
	}
}

func (m *MockPermissionRepository) CreatePermission(permission *model.Permission) error {
	if permission.Name == "" {
		return errors.New("nama permission tidak boleh kosong")
	}
	m.permissions[permission.ID] = permission
	return nil
}

func (m *MockPermissionRepository) GetPermissionByID(id string) (*model.Permission, error) {
	if permission, exists := m.permissions[id]; exists {
		return permission, nil
	}
	return nil, nil
}

func (m *MockPermissionRepository) GetPermissionByName(name string) (*model.Permission, error) {
	for _, permission := range m.permissions {
		if permission.Name == name {
			return permission, nil
		}
	}
	return nil, nil
}

func (m *MockPermissionRepository) GetAllPermissions() ([]*model.Permission, error) {
	var permissions []*model.Permission
	for _, permission := range m.permissions {
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (m *MockPermissionRepository) UpdatePermission(permission *model.Permission) error {
	if _, exists := m.permissions[permission.ID]; !exists {
		return errors.New("permission tidak ditemukan")
	}
	m.permissions[permission.ID] = permission
	return nil
}

func (m *MockPermissionRepository) DeletePermission(id string) error {
	if _, exists := m.permissions[id]; !exists {
		return errors.New("permission tidak ditemukan")
	}
	delete(m.permissions, id)
	return nil
}

func (m *MockPermissionRepository) GetPermissionsByRoleID(roleID string) ([]string, error) {
	// Mock implementation - return empty slice for now
	return []string{}, nil
}
