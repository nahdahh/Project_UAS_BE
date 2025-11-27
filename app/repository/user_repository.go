package repository

import (
	"database/sql"
	"uas_be/app/model"
)

// UserRepository adalah interface untuk akses data user dari database
type UserRepository interface {
	// CreateUser membuat user baru di database
	CreateUser(user *model.User) error
	
	// GetUserByID mengambil user berdasarkan ID
	GetUserByID(id string) (*model.User, error)
	
	// GetUserByUsername mengambil user berdasarkan username
	GetUserByUsername(username string) (*model.User, error)
	
	// GetUserByEmail mengambil user berdasarkan email
	GetUserByEmail(email string) (*model.User, error)
	
	// GetAllUsers mengambil semua user dengan pagination
	GetAllUsers(page, pageSize int) ([]*model.UserWithRole, int, error)
	
	// UpdateUser mengubah data user
	UpdateUser(user *model.User) error
	
	// DeleteUser menghapus user
	DeleteUser(id string) error
	
	// GetUserPermissions mengambil semua permission dari user berdasarkan role
	GetUserPermissions(userID string) ([]string, error)
}

// userRepositoryImpl adalah implementasi dari UserRepository
type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository membuat instance repository user baru
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// CreateUser membuat user baru
func (r *userRepositoryImpl) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive)
	return err
}

// GetUserByID mengambil user dari database berdasarkan ID
func (r *userRepositoryImpl) GetUserByID(id string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE id = $1`
	
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // user tidak ditemukan
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserByUsername mengambil user berdasarkan username
func (r *userRepositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE username = $1`
	
	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserByEmail mengambil user berdasarkan email
func (r *userRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE email = $1`
	
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return user, nil
}

// GetAllUsers mengambil semua user dengan pagination
func (r *userRepositoryImpl) GetAllUsers(page, pageSize int) ([]*model.UserWithRole, int, error) {
	// Hitung offset untuk pagination
	offset := (page - 1) * pageSize
	
	// Query untuk total items
	countQuery := `SELECT COUNT(*) FROM users`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}
	
	// Query untuk mengambil data user dengan role name dan permissions
	query := `
		SELECT u.id, u.username, u.email, u.full_name, r.name, u.is_active, u.created_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []*model.UserWithRole
	for rows.Next() {
		user := &model.UserWithRole{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.RoleName, &user.IsActive, &user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		
		// Ambil permissions untuk user ini
		permissions, err := r.GetUserPermissions(user.ID)
		if err == nil {
			user.Permissions = permissions
		}
		
		users = append(users, user)
	}
	
	return users, totalItems, nil
}

// UpdateUser mengubah data user
func (r *userRepositoryImpl) UpdateUser(user *model.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, password_hash = $3, full_name = $4, role_id = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
	`
	
	_, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive, user.ID)
	return err
}

// DeleteUser menghapus user dari database
func (r *userRepositoryImpl) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetUserPermissions mengambil semua permission dari user berdasarkan role
func (r *userRepositoryImpl) GetUserPermissions(userID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN users u ON rp.role_id = u.role_id
		WHERE u.id = $1
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	
	return permissions, nil
}
