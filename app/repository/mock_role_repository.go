package repository

import (
	"errors"
	"uas_be/app/model"
)

type MockRoleRepository struct {
	roles           map[string]*model.Role
	rolePermissions map[string][]string
}

func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles:           make(map[string]*model.Role),
		rolePermissions: make(map[string][]string),
	}
}

func (m *MockRoleRepository) CreateRole(role *model.Role) error {
	if role.Name == "" {
		return errors.New("nama role tidak boleh kosong")
	}
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) GetRoleByID(id string) (*model.Role, error) {
	if role, exists := m.roles[id]; exists {
		return role, nil
	}
	return nil, nil
}

func (m *MockRoleRepository) GetRoleByName(name string) (*model.Role, error) {
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil
}

func (m *MockRoleRepository) GetAllRoles() ([]*model.Role, error) {
	var roles []*model.Role
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (m *MockRoleRepository) UpdateRole(role *model.Role) error {
	if _, exists := m.roles[role.ID]; !exists {
		return errors.New("role tidak ditemukan")
	}
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) DeleteRole(id string) error {
	if _, exists := m.roles[id]; !exists {
		return errors.New("role tidak ditemukan")
	}
	delete(m.roles, id)
	return nil
}

func (m *MockRoleRepository) AssignPermissionToRole(roleID, permissionID string) error {
	if _, exists := m.roles[roleID]; !exists {
		return errors.New("role tidak ditemukan")
	}
	m.rolePermissions[roleID] = append(m.rolePermissions[roleID], permissionID)
	return nil
}

func (m *MockRoleRepository) RemovePermissionFromRole(roleID, permissionID string) error {
	if _, exists := m.roles[roleID]; !exists {
		return errors.New("role tidak ditemukan")
	}
	permissions := m.rolePermissions[roleID]
	for i, p := range permissions {
		if p == permissionID {
			m.rolePermissions[roleID] = append(permissions[:i], permissions[i+1:]...)
			break
		}
	}
	return nil
}

func (m *MockRoleRepository) GetRolePermissions(roleID string) ([]string, error) {
	if _, exists := m.roles[roleID]; !exists {
		return nil, errors.New("role tidak ditemukan")
	}
	return m.rolePermissions[roleID], nil
}
