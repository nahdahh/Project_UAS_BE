package repository

import (
	"errors"
	"uas_be/app/model"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	users map[string]*model.User
}

// NewMockUserRepository membuat instance mock repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	if user.Username == "" {
		return errors.New("username tidak boleh kosong")
	}
	if user.Email == "" {
		return errors.New("email tidak boleh kosong")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("id tidak boleh kosong")
	}
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*model.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
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

func (m *MockUserRepository) GetAllUsers(page, pageSize int) ([]*model.User, int, error) {
	var users []*model.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (m *MockUserRepository) UpdateUser(user *model.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user tidak ditemukan")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) DeleteUser(id string) error {
	if _, exists := m.users[id]; !exists {
		return errors.New("user tidak ditemukan")
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) AssignRole(userID, roleID string) error {
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user tidak ditemukan")
	}
	user.RoleID = roleID
	return nil
}

func (m *MockUserRepository) GetUserPermissions(userID string) ([]string, error) {
	if _, exists := m.users[userID]; !exists {
		return nil, errors.New("user tidak ditemukan")
	}
	return []string{"read", "write"}, nil
}
