package repository

import (
	"database/sql"
	"uas_be/app/model"
)

// PermissionRepository adalah interface untuk akses data permission dari database
type PermissionRepository interface {
	// GetPermissionByID mengambil permission berdasarkan ID
	GetPermissionByID(id string) (*model.Permission, error)

	// GetPermissionsByRoleID mengambil semua permission untuk sebuah role
	GetPermissionsByRoleID(roleID string) ([]string, error)

	// GetAllPermissions mengambil semua permission yang tersedia
	GetAllPermissions() ([]*model.Permission, error)

	// CreatePermission menambahkan permission baru ke database
	CreatePermission(permission *model.Permission) error
}

// permissionRepositoryImpl adalah implementasi dari PermissionRepository
type permissionRepositoryImpl struct {
	db *sql.DB
}

// NewPermissionRepository membuat instance repository permission baru
func NewPermissionRepository(db *sql.DB) PermissionRepository {
	return &permissionRepositoryImpl{db: db}
}

// GetPermissionsByRoleID mengambil semua permission untuk sebuah role berdasarkan role ID
func (r *permissionRepositoryImpl) GetPermissionsByRoleID(roleID string) ([]string, error) {
	query := `
		SELECT p.name FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := r.db.Query(query, roleID)
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

// GetAllPermissions mengambil semua permission yang tersedia
func (r *permissionRepositoryImpl) GetAllPermissions() ([]*model.Permission, error) {
	query := `SELECT id, name, resource, action, description FROM permissions ORDER BY resource, action`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*model.Permission
	for rows.Next() {
		p := &model.Permission{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

// GetPermissionByID mengambil permission berdasarkan ID
func (r *permissionRepositoryImpl) GetPermissionByID(id string) (*model.Permission, error) {
	query := `SELECT id, name, resource, action, description FROM permissions WHERE id = $1`

	p := &model.Permission{}
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return p, nil
}

// CreatePermission menambahkan permission baru ke database
func (r *permissionRepositoryImpl) CreatePermission(permission *model.Permission) error {
	query := `
		INSERT INTO permissions (id, name, resource, action, description)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, permission.ID, permission.Name, permission.Resource, permission.Action, permission.Description)
	return err
}
