package repository

import (
	"database/sql"
	"uas_be/app/model"
)

// RoleRepository adalah interface untuk akses data role dari database
type RoleRepository interface {
	// CreateRole membuat role baru
	CreateRole(role *model.Role) error

	// GetRoleByID mengambil role berdasarkan ID
	GetRoleByID(id string) (*model.Role, error)

	// GetRoleByName mengambil role berdasarkan nama
	GetRoleByName(name string) (*model.Role, error)

	// GetAllRoles mengambil semua role
	GetAllRoles() ([]*model.Role, error)

	// UpdateRole mengubah data role
	UpdateRole(role *model.Role) error

	// AssignPermissionToRole menambahkan permission ke role
	AssignPermissionToRole(roleID, permissionID string) error

	// RemovePermissionFromRole menghapus permission dari role
	RemovePermissionFromRole(roleID, permissionID string) error
}

// roleRepositoryImpl adalah implementasi dari RoleRepository
type roleRepositoryImpl struct {
	db *sql.DB
}

// NewRoleRepository membuat instance repository role baru
func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepositoryImpl{db: db}
}

// GetRoleByID mengambil role berdasarkan ID
func (r *roleRepositoryImpl) GetRoleByID(id string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`

	role := &model.Role{}
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// GetRoleByName mengambil role berdasarkan nama
func (r *roleRepositoryImpl) GetRoleByName(name string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE name = $1`

	role := &model.Role{}
	err := r.db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// GetAllRoles mengambil semua role
func (r *roleRepositoryImpl) GetAllRoles() ([]*model.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*model.Role
	for rows.Next() {
		role := &model.Role{}
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// CreateRole membuat role baru
func (r *roleRepositoryImpl) CreateRole(role *model.Role) error {
	query := `INSERT INTO roles (id, name, description, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(query, role.ID, role.Name, role.Description)
	return err
}

// UpdateRole mengubah data role
func (r *roleRepositoryImpl) UpdateRole(role *model.Role) error {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.Exec(query, role.Name, role.Description, role.ID)
	return err
}

// AssignPermissionToRole menambahkan permission ke role
func (r *roleRepositoryImpl) AssignPermissionToRole(roleID, permissionID string) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, roleID, permissionID)
	return err
}

// RemovePermissionFromRole menghapus permission dari role
func (r *roleRepositoryImpl) RemovePermissionFromRole(roleID, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err := r.db.Exec(query, roleID, permissionID)
	return err
}
