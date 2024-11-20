package service

import (
	"github.com/jmoiron/sqlx"
)

type RoleService struct {
	db *sqlx.DB
}

func NewRoleService(db *sqlx.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) CreateRole(name, description string) error {
	_, err := s.db.Exec(
		"INSERT INTO roles (name, description) VALUES (?, ?)",
		name, description,
	)
	return err
}

func (s *RoleService) AssignPermission(roleID, permissionID int) error {
	_, err := s.db.Exec(
		"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
		roleID, permissionID,
	)
	return err
}

func (s *RoleService) GetRoleByID(id int) (*Role, error) {
	var role Role
	err := s.db.Get(&role, "SELECT * FROM roles WHERE id = ?", id)
	return &role, err
}
