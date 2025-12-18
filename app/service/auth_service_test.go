package service

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"time"
	"uas_be/app/model"
	"uas_be/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ==================== MOCK REPOSITORIES ====================

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	users map[string]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(id string) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*model.User, error) {
	if user, exists := m.users[username]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) GetAllUsers(page, pageSize int) ([]*model.UserWithRole, int, error) {
	// Mock implementation - return empty slice for now
	return []*model.UserWithRole{}, 0, nil
}

func (m *MockUserRepository) UpdateUser(user *model.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) DeleteUser(id string) error {
	for username, user := range m.users {
		if user.ID == id {
			delete(m.users, username)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) AssignRole(userID string, roleID string) error {
	for _, user := range m.users {
		if user.ID == userID {
			user.RoleID = roleID
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) GetUserPermissions(userID string) ([]string, error) {
	// Mock implementation - return empty permissions
	return []string{}, nil
}

// MockPermissionRepository adalah mock untuk PermissionRepository
type MockPermissionRepository struct {
	permissions map[string][]string
}

func NewMockPermissionRepository() *MockPermissionRepository {
	return &MockPermissionRepository{
		permissions: make(map[string][]string),
	}
}

func (m *MockPermissionRepository) GetPermissionsByRoleID(roleID string) ([]string, error) {
	if perms, exists := m.permissions[roleID]; exists {
		return perms, nil
	}
	return []string{}, nil
}

func (m *MockPermissionRepository) GetPermissionByID(id string) (*model.Permission, error) {
	// Mock implementation - return a dummy permission
	return &model.Permission{
		ID:   id,
		Name: "test_permission",
	}, nil
}

func (m *MockPermissionRepository) CreatePermission(permission *model.Permission) error {
	return nil
}

func (m *MockPermissionRepository) GetAllPermissions() ([]*model.Permission, error) {
	return []*model.Permission{}, nil
}

func (m *MockPermissionRepository) DeletePermission(id string) error {
	return nil
}

// MockRoleRepository adalah mock untuk RoleRepository
type MockRoleRepository struct {
	roles map[string]*model.Role
}

func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles: make(map[string]*model.Role),
	}
}

func (m *MockRoleRepository) GetRoleByID(id string) (*model.Role, error) {
	if role, exists := m.roles[id]; exists {
		return role, nil
	}
	return nil, errors.New("role not found")
}

func (m *MockRoleRepository) GetRoleByName(name string) (*model.Role, error) {
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil
}

func (m *MockRoleRepository) CreateRole(role *model.Role) error {
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) GetAllRoles() ([]*model.Role, error) {
	roles := make([]*model.Role, 0, len(m.roles))
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (m *MockRoleRepository) UpdateRole(role *model.Role) error {
	if _, exists := m.roles[role.ID]; !exists {
		return errors.New("role not found")
	}
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) DeleteRole(id string) error {
	if _, exists := m.roles[id]; !exists {
		return errors.New("role not found")
	}
	delete(m.roles, id)
	return nil
}

func (m *MockRoleRepository) AssignPermissionToRole(roleID, permissionID string) error {
	// Mock implementation - just return nil
	return nil
}

func (m *MockRoleRepository) RemovePermissionFromRole(roleID, permissionID string) error {
	// Mock implementation - just return nil
	return nil
}

// ==================== TEST HELPER FUNCTIONS ====================

func setupAuthServiceTest() (AuthService, *MockUserRepository, *MockPermissionRepository, *MockRoleRepository) {
	// Initialize JWT for testing
	utils.InitJWT("test_secret_key_for_auth_service")

	userRepo := NewMockUserRepository()
	permRepo := NewMockPermissionRepository()
	roleRepo := NewMockRoleRepository()

	// Setup test data
	adminRoleID := "role_admin"
	roleRepo.CreateRole(&model.Role{
		ID:          adminRoleID,
		Name:        "Admin",
		Description: "Administrator role",
	})

	studentRoleID := "role_student"
	roleRepo.CreateRole(&model.Role{
		ID:          studentRoleID,
		Name:        "Mahasiswa",
		Description: "Student role",
	})

	permRepo.permissions[adminRoleID] = []string{"read", "write", "delete", "verify"}
	permRepo.permissions[studentRoleID] = []string{"read", "create"}

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	userRepo.CreateUser(&model.User{
		ID:           "user123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		RoleID:       adminRoleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	service := NewAuthService(userRepo, permRepo, roleRepo)
	return service, userRepo, permRepo, roleRepo
}

// ==================== UNIT TESTS ====================

// TestLoginSuccess menguji login dengan credentials yang benar
func TestLoginSuccess(t *testing.T) {
	// ARRANGE
	service, _, _, _ := setupAuthServiceTest()
	app := fiber.New()

	app.Post("/login", service.Login)

	// ACT
	body := `{"username":"testuser","password":"password123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// ASSERT
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	_ = resp // Avoid unused variable warning
}

// TestLoginInvalidCredentials menguji login dengan password salah
func TestLoginInvalidCredentials(t *testing.T) {
	// ARRANGE
	service, _, _, _ := setupAuthServiceTest()
	app := fiber.New()

	app.Post("/login", service.Login)

	tests := []struct {
		name     string
		username string
		password string
		wantCode int
	}{
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrongpassword",
			wantCode: fiber.StatusUnauthorized,
		},
		{
			name:     "non-existent user",
			username: "nonexistent",
			password: "password123",
			wantCode: fiber.StatusUnauthorized,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantCode: fiber.StatusBadRequest,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			wantCode: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			var body string
			switch tt.name {
			case "wrong password":
				body = `{"username":"testuser","password":"wrongpassword"}`
			case "non-existent user":
				body = `{"username":"nonexistent","password":"password123"}`
			case "empty username":
				body = `{"username":"","password":"password123"}`
			case "empty password":
				body = `{"username":"testuser","password":""}`
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			// ASSERT
			if tt.wantCode == fiber.StatusOK && resp.StatusCode != 200 {
				t.Errorf("Expected status 200 for valid credentials, got %d", resp.StatusCode)
			}
			if tt.wantCode != fiber.StatusOK && resp.StatusCode == 200 {
				t.Errorf("Expected error status for invalid credentials, got %d", resp.StatusCode)
			}
			_ = resp // Avoid unused variable warning
		})
	}
}

// TestLoginInactiveUser menguji login dengan user yang tidak aktif
func TestLoginInactiveUser(t *testing.T) {
	// ARRANGE
	service, userRepo, _, roleRepo := setupAuthServiceTest()

	// Create inactive user
	inactiveRoleID := "role_inactive"
	roleRepo.CreateRole(&model.Role{
		ID:   inactiveRoleID,
		Name: "InactiveRole",
	})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	userRepo.CreateUser(&model.User{
		ID:           "inactive_user",
		Username:     "inactiveuser",
		Email:        "inactive@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Inactive User",
		RoleID:       inactiveRoleID,
		IsActive:     false, // User is not active
	})

	app := fiber.New()
	app.Post("/login", service.Login)

	// ACT
	body := `{"username":"inactiveuser","password":"password123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// ASSERT
	if resp.StatusCode == 200 {
		t.Error("Login() should return error for inactive user")
	}
	_ = resp // Avoid unused variable warning
}

// TestLoginInvalidJSON menguji login dengan JSON yang tidak valid
func TestLoginInvalidJSON(t *testing.T) {
	// ARRANGE
	service, _, _, _ := setupAuthServiceTest()
	app := fiber.New()

	app.Post("/login", service.Login)

	// ACT
	body := `{invalid json`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// ASSERT
	if resp.StatusCode == 200 {
		t.Error("Login() should return error for invalid JSON")
	}
	_ = resp // Avoid unused variable warning
}
