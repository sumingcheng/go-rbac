package service

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Permission struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type PermissionService struct {
	db *sqlx.DB
}

func NewPermissionService(db *sqlx.DB) *PermissionService {
	return &PermissionService{db: db}
}

func (s *PermissionService) CreatePermission(name, description string) error {
	_, err := s.db.Exec(
		"INSERT INTO permissions (name, description) VALUES (?, ?)",
		name, description,
	)
	return err
}

func (s *PermissionService) GetPermissionByID(id int) (*Permission, error) {
	var permission Permission
	err := s.db.Get(&permission, "SELECT * FROM permissions WHERE id = ?", id)
	return &permission, err
}
