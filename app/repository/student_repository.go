package repository

import (
	"database/sql"
	"uas_be/app/model"
)

// StudentRepository adalah interface untuk akses data student dari database
type StudentRepository interface {
	// CreateStudent membuat student baru
	CreateStudent(student *model.Student) error
	
	// GetStudentByID mengambil student berdasarkan ID
	GetStudentByID(id string) (*model.Student, error)
	
	// GetStudentByUserID mengambil student berdasarkan user ID
	GetStudentByUserID(userID string) (*model.Student, error)
	
	// GetStudentByStudentID mengambil student berdasarkan NIM
	GetStudentByStudentID(studentID string) (*model.Student, error)
	
	// GetAllStudents mengambil semua student dengan pagination
	GetAllStudents(page, pageSize int) ([]*model.StudentWithUser, int, error)
	
	// GetStudentsByAdvisorID mengambil student berdasarkan dosen wali
	GetStudentsByAdvisorID(advisorID string) ([]*model.StudentWithUser, error)
	
	// UpdateStudent mengubah data student
	UpdateStudent(student *model.Student) error
	
	// DeleteStudent menghapus student
	DeleteStudent(id string) error
}

// studentRepositoryImpl adalah implementasi dari StudentRepository
type studentRepositoryImpl struct {
	db *sql.DB
}

// NewStudentRepository membuat instance repository student baru
func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepositoryImpl{db: db}
}

// CreateStudent membuat student baru di database
func (r *studentRepositoryImpl) CreateStudent(student *model.Student) error {
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Exec(query, student.ID, student.UserID, student.StudentID, 
		student.ProgramStudy, student.AcademicYear, student.AdvisorID)
	return err
}

// GetStudentByID mengambil student berdasarkan ID
func (r *studentRepositoryImpl) GetStudentByID(id string) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
		FROM students WHERE id = $1
	`
	
	student := &model.Student{}
	err := r.db.QueryRow(query, id).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID, &student.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return student, nil
}

// GetStudentByUserID mengambil student berdasarkan user ID
func (r *studentRepositoryImpl) GetStudentByUserID(userID string) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
		FROM students WHERE user_id = $1
	`
	
	student := &model.Student{}
	err := r.db.QueryRow(query, userID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID, &student.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return student, nil
}

// GetStudentByStudentID mengambil student berdasarkan NIM
func (r *studentRepositoryImpl) GetStudentByStudentID(studentID string) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
		FROM students WHERE student_id = $1
	`
	
	student := &model.Student{}
	err := r.db.QueryRow(query, studentID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID, &student.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return student, nil
}

// GetAllStudents mengambil semua student dengan pagination
func (r *studentRepositoryImpl) GetAllStudents(page, pageSize int) ([]*model.StudentWithUser, int, error) {
	offset := (page - 1) * pageSize
	
	// Hitung total items
	countQuery := `SELECT COUNT(*) FROM students`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}
	
	// Query untuk mengambil data student dengan user info
	query := `
		SELECT s.id, s.user_id, s.student_id, u.full_name, u.email, 
		       s.program_study, s.academic_year, s.advisor_id, s.created_at
		FROM students s
		JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var students []*model.StudentWithUser
	for rows.Next() {
		student := &model.StudentWithUser{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentID, &student.FullName,
			&student.Email, &student.ProgramStudy, &student.AcademicYear, 
			&student.AdvisorID, &student.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		
		// Ambil nama advisor jika ada
		if student.AdvisorID != "" {
			advisorQuery := `SELECT u.full_name FROM lecturers l JOIN users u ON l.user_id = u.id WHERE l.id = $1`
			r.db.QueryRow(advisorQuery, student.AdvisorID).Scan(&student.AdvisorName)
		}
		
		students = append(students, student)
	}
	
	return students, totalItems, nil
}

// GetStudentsByAdvisorID mengambil student berdasarkan dosen wali
func (r *studentRepositoryImpl) GetStudentsByAdvisorID(advisorID string) ([]*model.StudentWithUser, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, u.full_name, u.email,
		       s.program_study, s.academic_year, s.advisor_id, s.created_at
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.advisor_id = $1
		ORDER BY s.created_at DESC
	`
	
	rows, err := r.db.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var students []*model.StudentWithUser
	for rows.Next() {
		student := &model.StudentWithUser{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentID, &student.FullName,
			&student.Email, &student.ProgramStudy, &student.AcademicYear,
			&student.AdvisorID, &student.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	
	return students, nil
}

// UpdateStudent mengubah data student
func (r *studentRepositoryImpl) UpdateStudent(student *model.Student) error {
	query := `
		UPDATE students 
		SET program_study = $1, academic_year = $2, advisor_id = $3
		WHERE id = $4
	`
	
	_, err := r.db.Exec(query, student.ProgramStudy, student.AcademicYear, student.AdvisorID, student.ID)
	return err
}

// DeleteStudent menghapus student
func (r *studentRepositoryImpl) DeleteStudent(id string) error {
	query := `DELETE FROM students WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
