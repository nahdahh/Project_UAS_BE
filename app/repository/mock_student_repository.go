package repository

import (
	"errors"
	"uas_be/app/model"
)

type MockStudentRepository struct {
	students map[string]*model.Student
}

func NewMockStudentRepository() *MockStudentRepository {
	return &MockStudentRepository{
		students: make(map[string]*model.Student),
	}
}

func (m *MockStudentRepository) CreateStudent(student *model.Student) error {
	if student.StudentID == "" {
		return errors.New("student ID tidak boleh kosong")
	}
	m.students[student.ID] = student
	return nil
}

func (m *MockStudentRepository) GetStudentByID(id string) (*model.Student, error) {
	if student, exists := m.students[id]; exists {
		return student, nil
	}
	return nil, nil
}

func (m *MockStudentRepository) GetStudentByUserID(userID string) (*model.Student, error) {
	for _, student := range m.students {
		if student.UserID == userID {
			return student, nil
		}
	}
	return nil, nil
}

func (m *MockStudentRepository) GetStudentByStudentID(studentID string) (*model.Student, error) {
	for _, student := range m.students {
		if student.StudentID == studentID {
			return student, nil
		}
	}
	return nil, nil
}

func (m *MockStudentRepository) GetAllStudents(page, pageSize int) ([]*model.StudentWithUser, int, error) {
	var students []*model.StudentWithUser
	for _, student := range m.students {
		studentWithUser := &model.StudentWithUser{
			ID:           student.ID,
			UserID:       student.UserID,
			StudentID:    student.StudentID,
			ProgramStudy: student.ProgramStudy,
			AcademicYear: student.AcademicYear,
			AdvisorID:    student.AdvisorID,
			CreatedAt:    student.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		students = append(students, studentWithUser)
	}
	return students, len(students), nil
}

func (m *MockStudentRepository) GetStudentsByAdvisorID(advisorID string) ([]*model.StudentWithUser, error) {
	var students []*model.StudentWithUser
	for _, student := range m.students {
		if student.AdvisorID == advisorID {
			studentWithUser := &model.StudentWithUser{
				ID:           student.ID,
				UserID:       student.UserID,
				StudentID:    student.StudentID,
				ProgramStudy: student.ProgramStudy,
				AcademicYear: student.AcademicYear,
				AdvisorID:    student.AdvisorID,
				CreatedAt:    student.CreatedAt.Format("2006-01-02T15:04:05Z"),
			}
			students = append(students, studentWithUser)
		}
	}
	return students, nil
}

func (m *MockStudentRepository) UpdateStudent(student *model.Student) error {
	if _, exists := m.students[student.ID]; !exists {
		return errors.New("student tidak ditemukan")
	}
	m.students[student.ID] = student
	return nil
}

func (m *MockStudentRepository) DeleteStudent(id string) error {
	if _, exists := m.students[id]; !exists {
		return errors.New("student tidak ditemukan")
	}
	delete(m.students, id)
	return nil
}

func (m *MockStudentRepository) UpdateAdvisor(studentID, advisorID string) error {
	student, exists := m.students[studentID]
	if !exists {
		return errors.New("student tidak ditemukan")
	}
	student.AdvisorID = advisorID
	return nil
}
