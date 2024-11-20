package database

import (
	"github.com/jmoiron/sqlx"
)

// InitBaseData 初始化基础数据
func InitBaseData(db *sqlx.DB) error {
	// 创建超级管理员角色
	if err := createSuperAdminRole(db); err != nil {
		return err
	}

	// 创建基础权限
	if err := createBasePermissions(db); err != nil {
		return err
	}

	// 为超级管理员分配所有权限
	return assignSuperAdminPermissions(db)
}

func createSuperAdminRole(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT IGNORE INTO roles (name, description) 
		VALUES ('super_admin', '超级管理员');
	`)
	return err
}

func createBasePermissions(db *sqlx.DB) error {
	basePermissions := []struct {
		Name        string
		Description string
		Resource    string
		Action      string
	}{
		{"user_create", "创建用户", "user", "create"},
		{"user_read", "读取用户信息", "user", "read"},
		{"user_update", "更新用户信息", "user", "update"},
		{"user_delete", "删除用户", "user", "delete"},
		{"role_create", "创建角色", "role", "create"},
		{"role_read", "读取角色信息", "role", "read"},
		{"role_update", "更新角色信息", "role", "update"},
		{"role_delete", "删除角色", "role", "delete"},
		{"permission_create", "创建权限", "permission", "create"},
		{"permission_read", "读取权限信息", "permission", "read"},
		{"permission_update", "更新权限信息", "permission", "update"},
		{"permission_delete", "删除权限", "permission", "delete"},
	}

	for _, perm := range basePermissions {
		_, err := db.Exec(`
			INSERT IGNORE INTO permissions (name, description, resource, action) 
			VALUES (?, ?, ?, ?)
		`, perm.Name, perm.Description, perm.Resource, perm.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func assignSuperAdminPermissions(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT IGNORE INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id
		FROM roles r
		CROSS JOIN permissions p
		WHERE r.name = 'super_admin'
	`)
	return err
}
