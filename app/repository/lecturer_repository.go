package repository

import (
	"database/sql"
	"uas_be/app/model"
)

// LecturerRepository adalah interface untuk akses data lecturer dari database
type LecturerRepository interface {
	// CreateLecturer membuat lecturer baru
	CreateLecturer(lecturer *model.Lecturer) error
	
	// GetLecturerByID mengambil lecturer berdasarkan ID
	GetLecturerByID(id string) (*model.Lecturer, error)
	
	// GetLecturerByUserID mengambil lecturer berdasarkan user ID
	GetLecturerByUserID(userID string) (*model.Lecturer, error)
	
	// GetLecturerByLecturerID mengambil lecturer berdasarkan NIDN
	GetLecturerByLecturerID(lecturerID string) (*model.Lecturer, error)
	
	// GetAllLecturers mengambil semua lecturer dengan pagination
	GetAllLecturers(page, pageSize int) ([]*model.LecturerWithUser, int, error)
	
	// UpdateLecturer mengubah data lecturer
	UpdateLecturer(lecturer *model.Lecturer) error
	
	// DeleteLecturer menghapus lecturer
	DeleteLecturer(id string) error
}

// lecturerRepositoryImpl adalah implementasi dari LecturerRepository
type lecturerRepositoryImpl struct {
	db *sql.DB
}

// NewLecturerRepository membuat instance repository lecturer baru
func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepositoryImpl{db: db}
}

// CreateLecturer membuat lecturer baru di database
func (r *lecturerRepositoryImpl) CreateLecturer(lecturer *model.Lecturer) error {
	query := `
		INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := r.db.Exec(query, lecturer.ID, lecturer.UserID, lecturer.LecturerID, lecturer.Department)
	return err
}

// GetLecturerByID mengambil lecturer berdasarkan ID
func (r *lecturerRepositoryImpl) GetLecturerByID(id string) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers WHERE id = $1
	`
	
	lecturer := &model.Lecturer{}
	err := r.db.QueryRow(query, id).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.Department, &lecturer.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return lecturer, nil
}

// GetLecturerByUserID mengambil lecturer berdasarkan user ID
func (r *lecturerRepositoryImpl) GetLecturerByUserID(userID string) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers WHERE user_id = $1
	`
	
	lecturer := &model.Lecturer{}
	err := r.db.QueryRow(query, userID).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.Department, &lecturer.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return lecturer, nil
}

// GetLecturerByLecturerID mengambil lecturer berdasarkan NIDN
func (r *lecturerRepositoryImpl) GetLecturerByLecturerID(lecturerID string) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers WHERE lecturer_id = $1
	`
	
	lecturer := &model.Lecturer{}
	err := r.db.QueryRow(query, lecturerID).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.Department, &lecturer.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return lecturer, nil
}

// GetAllLecturers mengambil semua lecturer dengan pagination
func (r *lecturerRepositoryImpl) GetAllLecturers(page, pageSize int) ([]*model.LecturerWithUser, int, error) {
	offset := (page - 1) * pageSize
	
	// Hitung total items
	countQuery := `SELECT COUNT(*) FROM lecturers`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}
	
	// Query untuk mengambil data lecturer dengan user info
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, u.full_name, u.email, l.department, l.created_at
		FROM lecturers l
		JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var lecturers []*model.LecturerWithUser
	for rows.Next() {
		lecturer := &model.LecturerWithUser{}
		err := rows.Scan(
			&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.FullName,
			&lecturer.Email, &lecturer.Department, &lecturer.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		lecturers = append(lecturers, lecturer)
	}
	
	return lecturers, totalItems, nil
}

// UpdateLecturer mengubah data lecturer
func (r *lecturerRepositoryImpl) UpdateLecturer(lecturer *model.Lecturer) error {
	query := `
		UPDATE lecturers 
		SET department = $1
		WHERE id = $2
	`
	
	_, err := r.db.Exec(query, lecturer.Department, lecturer.ID)
	return err
}

// DeleteLecturer menghapus lecturer
func (r *lecturerRepositoryImpl) DeleteLecturer(id string) error {
	query := `DELETE FROM lecturers WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
