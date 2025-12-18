package service

import (
	"testing"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/google/uuid"
)

// TestCreateLecturer menguji pembuatan lecturer baru
func TestCreateLecturer(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockLecturerRepository()

	tests := []struct {
		name       string
		userID     string
		lecturerID string
		department string
		wantErr    bool
	}{
		{
			name:       "Valid Lecturer",
			userID:     uuid.New().String(),
			lecturerID: "198801012015041001",
			department: "Teknik Informatika",
			wantErr:    false,
		},
		{
			name:       "Empty Lecturer ID",
			userID:     uuid.New().String(),
			lecturerID: "",
			department: "Teknik Informatika",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			lecturer := &model.Lecturer{
				ID:         uuid.New().String(),
				UserID:     tt.userID,
				LecturerID: tt.lecturerID,
				Department: tt.department,
			}

			err := mockRepo.CreateLecturer(lecturer)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLecturer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetLecturerByID menguji pengambilan lecturer berdasarkan ID
func TestGetLecturerByID(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockLecturerRepository()
	lecturerID := uuid.New().String()
	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     uuid.New().String(),
		LecturerID: "198801012015041001",
		Department: "Teknik Informatika",
	}
	mockRepo.CreateLecturer(lecturer)

	tests := []struct {
		name    string
		id      string
		wantNil bool
	}{
		{"Valid ID", lecturerID, false},
		{"Invalid ID", uuid.New().String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := mockRepo.GetLecturerByID(tt.id)

			// ASSERT
			if err != nil {
				t.Errorf("GetLecturerByID() unexpected error = %v", err)
			}
			if tt.wantNil && result != nil {
				t.Errorf("expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("expected lecturer, got nil")
			}
		})
	}
}

// TestGetLecturerByUserID menguji pengambilan lecturer berdasarkan user ID
func TestGetLecturerByUserID(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockLecturerRepository()
	userID := uuid.New().String()
	lecturer := &model.Lecturer{
		ID:         uuid.New().String(),
		UserID:     userID,
		LecturerID: "198801012015041001",
		Department: "Teknik Informatika",
	}
	mockRepo.CreateLecturer(lecturer)

	tests := []struct {
		name    string
		userID  string
		wantNil bool
	}{
		{"Valid User ID", userID, false},
		{"Invalid User ID", uuid.New().String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := mockRepo.GetLecturerByUserID(tt.userID)

			// ASSERT
			if err != nil {
				t.Errorf("GetLecturerByUserID() unexpected error = %v", err)
			}
			if tt.wantNil && result != nil {
				t.Errorf("expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("expected lecturer, got nil")
			}
		})
	}
}

// TestUpdateLecturer menguji update lecturer
func TestUpdateLecturer(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockLecturerRepository()
	lecturerID := uuid.New().String()
	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     uuid.New().String(),
		LecturerID: "198801012015041001",
		Department: "Teknik Informatika",
	}
	mockRepo.CreateLecturer(lecturer)

	tests := []struct {
		name          string
		id            string
		newDepartment string
		wantErr       bool
	}{
		{"Valid Update", lecturerID, "Sistem Informasi", false},
		{"Invalid ID", uuid.New().String(), "Department Baru", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			updatedLecturer := &model.Lecturer{
				ID:         tt.id,
				Department: tt.newDepartment,
			}
			err := mockRepo.UpdateLecturer(updatedLecturer)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLecturer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDeleteLecturer menguji penghapusan lecturer
func TestDeleteLecturer(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockLecturerRepository()
	lecturerID := uuid.New().String()
	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     uuid.New().String(),
		LecturerID: "198801012015041001",
	}
	mockRepo.CreateLecturer(lecturer)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"Valid Delete", lecturerID, false},
		{"Already Deleted", lecturerID, true},
		{"Invalid ID", uuid.New().String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			err := mockRepo.DeleteLecturer(tt.id)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLecturer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
