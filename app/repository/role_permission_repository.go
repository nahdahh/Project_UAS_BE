package repository

import (
	"database/sql"
	"uas_be/app/model"
)

type RolePermissionRepository struct {
	db *sql.DB
}

func NewRolePermissionRepository(db *sql.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

// Tambahkan permission ke role
func (r *RolePermissionRepository) AssignPermission(roleID, permissionID string) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(query, roleID, permissionID)
	return err
}

// Hapus permission dari role
func (r *RolePermissionRepository) RemovePermission(roleID, permissionID string) error {
	_, err := r.db.Exec(`DELETE FROM role_permissions WHERE role_id=$1 AND permission_id=$2`,
		roleID, permissionID)
	return err
}

// Ambil semua permission dari role
func (r *RolePermissionRepository) GetPermissionsByRole(roleID string) ([]model.Permission, error) {
	query := `
		SELECT p.id, p.name, p.resource, p.action, p.description
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var perm model.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}
