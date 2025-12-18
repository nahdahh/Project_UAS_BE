package service

import (
	"testing"
	"uas_be/app/model"
	"uas_be/app/repository"

	"github.com/google/uuid"
)

// TestCreateStudent menguji pembuatan student baru
func TestCreateStudent(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockStudentRepository()

	tests := []struct {
		name         string
		userID       string
		studentID    string
		programStudy string
		academicYear string
		advisorID    string
		wantErr      bool
	}{
		{
			name:         "Valid Student",
			userID:       uuid.New().String(),
			studentID:    "162211001",
			programStudy: "Teknik Informatika",
			academicYear: "2022",
			advisorID:    uuid.New().String(),
			wantErr:      false,
		},
		{
			name:         "Empty Student ID",
			userID:       uuid.New().String(),
			studentID:    "",
			programStudy: "Teknik Informatika",
			academicYear: "2022",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			student := &model.Student{
				ID:           uuid.New().String(),
				UserID:       tt.userID,
				StudentID:    tt.studentID,
				ProgramStudy: tt.programStudy,
				AcademicYear: tt.academicYear,
				AdvisorID:    tt.advisorID,
			}

			err := mockRepo.CreateStudent(student)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetStudentByID menguji pengambilan student berdasarkan ID
func TestGetStudentByID(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockStudentRepository()
	studentID := uuid.New().String()
	student := &model.Student{
		ID:        studentID,
		UserID:    uuid.New().String(),
		StudentID: "162211001",
	}
	mockRepo.CreateStudent(student)

	tests := []struct {
		name    string
		id      string
		wantNil bool
	}{
		{"Valid ID", studentID, false},
		{"Invalid ID", uuid.New().String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result, err := mockRepo.GetStudentByID(tt.id)

			// ASSERT
			if err != nil {
				t.Errorf("GetStudentByID() unexpected error = %v", err)
			}
			if tt.wantNil && result != nil {
				t.Errorf("expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("expected student, got nil")
			}
		})
	}
}

// TestGetStudentsByAdvisorID menguji pengambilan students berdasarkan advisor
func TestGetStudentsByAdvisorID(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockStudentRepository()
	advisorID := uuid.New().String()

	// Buat beberapa students dengan advisor yang sama
	for i := 1; i <= 3; i++ {
		student := &model.Student{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			StudentID: "16221100" + string(rune(i)),
			AdvisorID: advisorID,
		}
		mockRepo.CreateStudent(student)
	}

	tests := []struct {
		name        string
		advisorID   string
		expectedLen int
	}{
		{"Valid Advisor with Students", advisorID, 3},
		{"Invalid Advisor", uuid.New().String(), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			students, err := mockRepo.GetStudentsByAdvisorID(tt.advisorID)

			// ASSERT
			if err != nil {
				t.Errorf("GetStudentsByAdvisorID() error = %v", err)
			}
			if len(students) != tt.expectedLen {
				t.Errorf("expected %d students, got %d", tt.expectedLen, len(students))
			}
		})
	}
}

// TestUpdateStudent menguji update student
func TestUpdateStudent(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockStudentRepository()
	studentID := uuid.New().String()
	student := &model.Student{
		ID:           studentID,
		UserID:       uuid.New().String(),
		StudentID:    "162211001",
		ProgramStudy: "Teknik Informatika",
	}
	mockRepo.CreateStudent(student)

	tests := []struct {
		name            string
		id              string
		newProgramStudy string
		wantErr         bool
	}{
		{"Valid Update", studentID, "Sistem Informasi", false},
		{"Invalid ID", uuid.New().String(), "Program Baru", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			updatedStudent := &model.Student{
				ID:           tt.id,
				ProgramStudy: tt.newProgramStudy,
			}
			err := mockRepo.UpdateStudent(updatedStudent)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDeleteStudent menguji penghapusan student
func TestDeleteStudent(t *testing.T) {
	// ARRANGE
	mockRepo := repository.NewMockStudentRepository()
	studentID := uuid.New().String()
	student := &model.Student{
		ID:        studentID,
		UserID:    uuid.New().String(),
		StudentID: "162211001",
	}
	mockRepo.CreateStudent(student)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"Valid Delete", studentID, false},
		{"Already Deleted", studentID, true},
		{"Invalid ID", uuid.New().String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			err := mockRepo.DeleteStudent(tt.id)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
