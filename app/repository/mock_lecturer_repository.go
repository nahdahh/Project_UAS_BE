package repository

import (
	"errors"
	"uas_be/app/model"
)

type MockLecturerRepository struct {
	lecturers map[string]*model.Lecturer
}

func NewMockLecturerRepository() *MockLecturerRepository {
	return &MockLecturerRepository{
		lecturers: make(map[string]*model.Lecturer),
	}
}

func (m *MockLecturerRepository) CreateLecturer(lecturer *model.Lecturer) error {
	if lecturer.LecturerID == "" {
		return errors.New("lecturer ID tidak boleh kosong")
	}
	m.lecturers[lecturer.ID] = lecturer
	return nil
}

func (m *MockLecturerRepository) GetLecturerByID(id string) (*model.Lecturer, error) {
	if lecturer, exists := m.lecturers[id]; exists {
		return lecturer, nil
	}
	return nil, nil
}

func (m *MockLecturerRepository) GetLecturerByUserID(userID string) (*model.Lecturer, error) {
	for _, lecturer := range m.lecturers {
		if lecturer.UserID == userID {
			return lecturer, nil
		}
	}
	return nil, nil
}

func (m *MockLecturerRepository) GetLecturerByLecturerID(lecturerID string) (*model.Lecturer, error) {
	for _, lecturer := range m.lecturers {
		if lecturer.LecturerID == lecturerID {
			return lecturer, nil
		}
	}
	return nil, nil
}

func (m *MockLecturerRepository) GetAllLecturers(page, pageSize int) ([]*model.LecturerWithUser, int, error) {
	var lecturers []*model.LecturerWithUser
	for _, lecturer := range m.lecturers {
		lecturerWithUser := &model.LecturerWithUser{
			ID:         lecturer.ID,
			UserID:     lecturer.UserID,
			LecturerID: lecturer.LecturerID,
			Department: lecturer.Department,
			CreatedAt:  lecturer.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		lecturers = append(lecturers, lecturerWithUser)
	}
	return lecturers, len(lecturers), nil
}

func (m *MockLecturerRepository) UpdateLecturer(lecturer *model.Lecturer) error {
	if _, exists := m.lecturers[lecturer.ID]; !exists {
		return errors.New("lecturer tidak ditemukan")
	}
	m.lecturers[lecturer.ID] = lecturer
	return nil
}

func (m *MockLecturerRepository) DeleteLecturer(id string) error {
	if _, exists := m.lecturers[id]; !exists {
		return errors.New("lecturer tidak ditemukan")
	}
	delete(m.lecturers, id)
	return nil
}
