package service

import (
	"testing"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TestCreateUser menguji pembuatan user baru
func TestCreateUser(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockUserRepository()

	tests := []struct {
		name     string
		username string
		email    string
		password string
		fullName string
		roleID   string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid User",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			fullName: "Test User",
			roleID:   uuid.New().String(),
			wantErr:  false,
		},
		{
			name:     "Empty Username",
			username: "",
			email:    "test2@example.com",
			password: "password123",
			fullName: "Test User 2",
			roleID:   uuid.New().String(),
			wantErr:  true,
			errMsg:   "username tidak boleh kosong",
		},
		{
			name:     "Empty Email",
			username: "testuser3",
			email:    "",
			password: "password123",
			fullName: "Test User 3",
			roleID:   uuid.New().String(),
			wantErr:  true,
			errMsg:   "email tidak boleh kosong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(tt.password), bcrypt.DefaultCost)
			user := &model.User{
				ID:           uuid.New().String(),
				Username:     tt.username,
				Email:        tt.email,
				PasswordHash: string(hashedPassword),
				FullName:     tt.fullName,
				RoleID:       tt.roleID,
				IsActive:     true,
			}

			var err error
			if tt.username != "" && tt.email != "" {
				err = mockRepo.CreateUser(user)
			} else {
				err = mockRepo.CreateUser(user)
			}

			// ASSERT
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

// TestGetUserByID menguji pengambilan user berdasarkan ID
func TestGetUserByID(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockUserRepository()
	userID := uuid.New().String()
	user := &model.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		IsActive: true,
	}
	mockRepo.CreateUser(user)

	tests := []struct {
		name    string
		id      string
		wantErr bool
		wantNil bool
	}{
		{"Valid ID", userID, false, false},
		{"Invalid ID", uuid.New().String(), false, true},
		{"Empty ID", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := mockRepo.GetUserByID(tt.id)

			// ASSERT
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if tt.wantNil && result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
			if !tt.wantNil && !tt.wantErr && result == nil {
				t.Errorf("expected user, got nil")
			}
		})
	}
}

// TestUpdateUser menguji update user
func TestUpdateUser(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockUserRepository()
	userID := uuid.New().String()
	user := &model.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.CreateUser(user)

	tests := []struct {
		name        string
		id          string
		newUsername string
		newEmail    string
		wantErr     bool
	}{
		{"Valid Update", userID, "updateduser", "updated@example.com", false},
		{"Invalid ID", uuid.New().String(), "user2", "user2@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			updatedUser := &model.User{
				ID:       tt.id,
				Username: tt.newUsername,
				Email:    tt.newEmail,
				FullName: "Updated User",
			}
			err := mockRepo.UpdateUser(updatedUser)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDeleteUser menguji penghapusan user
func TestDeleteUser(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockUserRepository()
	userID := uuid.New().String()
	user := &model.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}
	mockRepo.CreateUser(user)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"Valid Delete", userID, false},
		{"Already Deleted", userID, true},
		{"Invalid ID", uuid.New().String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			err := mockRepo.DeleteUser(tt.id)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
