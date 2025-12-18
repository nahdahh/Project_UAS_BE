package service

import (
	"testing"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/google/uuid"
)

// TestCreateRole menguji pembuatan role baru
func TestCreateRole(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	tests := []struct {
		name        string
		roleName    string
		description string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid Role",
			roleName:    "Admin",
			description: "Administrator role",
			wantErr:     false,
		},
		{
			name:        "Empty Name",
			roleName:    "",
			description: "Empty name role",
			wantErr:     true,
			errMsg:      "nama role tidak boleh kosong",
		},
		{
			name:        "Duplicate Name",
			roleName:    "Admin",
			description: "Duplicate admin",
			wantErr:     true,
			errMsg:      "nama role sudah terdaftar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			role, err := service.createRole(tt.roleName, tt.description)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("createRole() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("error message = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr && role == nil {
				t.Errorf("expected role, got nil")
			}
		})
	}
}

// TestGetRoleByID menguji pengambilan role berdasarkan ID
func TestGetRoleByID(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	roleID := uuid.New().String()
	role := &model.Role{
		ID:          roleID,
		Name:        "Admin",
		Description: "Administrator",
	}
	mockRoleRepo.CreateRole(role)

	tests := []struct {
		name    string
		id      string
		wantErr bool
		wantNil bool
	}{
		{"Valid ID", roleID, false, false},
		{"Invalid ID", uuid.New().String(), false, true},
		{"Empty ID", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := service.getRoleByID(tt.id)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("getRoleByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && result != nil {
				t.Errorf("expected nil, got %v", result)
			}
			if !tt.wantNil && !tt.wantErr && result == nil {
				t.Errorf("expected role, got nil")
			}
		})
	}
}

// TestGetRoleByName menguji pengambilan role berdasarkan nama
func TestGetRoleByName(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	role := &model.Role{
		ID:          uuid.New().String(),
		Name:        "Admin",
		Description: "Administrator",
	}
	mockRoleRepo.CreateRole(role)

	tests := []struct {
		name     string
		roleName string
		wantErr  bool
		wantNil  bool
	}{
		{"Valid Name", "Admin", false, false},
		{"Invalid Name", "NonExistent", false, true},
		{"Empty Name", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := service.getRoleByName(tt.roleName)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("getRoleByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && result != nil {
				t.Errorf("expected nil, got %v", result)
			}
			if !tt.wantNil && !tt.wantErr && result == nil {
				t.Errorf("expected role, got nil")
			}
		})
	}
}

// TestGetAllRoles menguji pengambilan semua roles
func TestGetAllRoles(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	// Buat beberapa roles
	roles := []*model.Role{
		{ID: uuid.New().String(), Name: "Admin", Description: "Administrator"},
		{ID: uuid.New().String(), Name: "User", Description: "Regular User"},
		{ID: uuid.New().String(), Name: "Moderator", Description: "Moderator"},
	}
	for _, role := range roles {
		mockRoleRepo.CreateRole(role)
	}

	tests := []struct {
		name        string
		expectedLen int
		wantErr     bool
	}{
		{"Get All Roles", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := service.getAllRoles()

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(result) != tt.expectedLen {
				t.Errorf("expected %d roles, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

// TestUpdateRole menguji update role
func TestUpdateRole(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	roleID := uuid.New().String()
	role := &model.Role{
		ID:          roleID,
		Name:        "Admin",
		Description: "Old Description",
	}
	mockRoleRepo.CreateRole(role)

	tests := []struct {
		name           string
		id             string
		newDescription string
		wantErr        bool
	}{
		{"Valid Update", roleID, "New Description", false},
		{"Invalid ID", uuid.New().String(), "Description", true},
		{"Empty ID", "", "Description", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			_, err := service.updateRole(tt.id, tt.newDescription)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("updateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAssignPermissionToRole menguji assign permission ke role
func TestAssignPermissionToRole(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	roleID := uuid.New().String()
	permissionID := uuid.New().String()

	role := &model.Role{ID: roleID, Name: "Admin", Description: "Admin Role"}
	permission := &model.Permission{ID: permissionID, Name: "create_user", Description: "Create User"}

	mockRoleRepo.CreateRole(role)
	mockPermissionRepo.CreatePermission(permission)

	tests := []struct {
		name         string
		roleID       string
		permissionID string
		wantErr      bool
	}{
		{"Valid Assignment", roleID, permissionID, false},
		{"Invalid Role ID", uuid.New().String(), permissionID, true},
		{"Invalid Permission ID", roleID, uuid.New().String(), true},
		{"Empty IDs", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			err := service.assignPermissionToRole(tt.roleID, tt.permissionID)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("assignPermissionToRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRemovePermissionFromRole menguji remove permission dari role
func TestRemovePermissionFromRole(t *testing.T) {
	// ARRANGE
	mockRoleRepo := repository.NewMockRoleRepository()
	mockPermissionRepo := repository.NewMockPermissionRepository()
	service := &roleServiceImpl{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
	}

	roleID := uuid.New().String()
	permissionID := uuid.New().String()

	role := &model.Role{ID: roleID, Name: "Admin", Description: "Admin Role"}
	mockRoleRepo.CreateRole(role)
	mockRoleRepo.AssignPermissionToRole(roleID, permissionID)

	tests := []struct {
		name         string
		roleID       string
		permissionID string
		wantErr      bool
	}{
		{"Valid Removal", roleID, permissionID, false},
		{"Invalid Role ID", uuid.New().String(), permissionID, true},
		{"Empty IDs", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			err := service.removePermissionFromRole(tt.roleID, tt.permissionID)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("removePermissionFromRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
