package database

import (
	"github.com/jmoiron/sqlx"
)

// InitSchema 初始化数据库表结构
func InitSchema(db *sqlx.DB) error {
	// 执行建表语句
	_, err := db.Exec(TableSchema)
	if err != nil {
		return err
	}

	// 初始化基础数据
	return InitBaseData(db)
}

// TableSchema 定义所有表结构
const TableSchema = `
	-- 用户表
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		status TINYINT DEFAULT 1 COMMENT '1:active, 0:inactive',
		token VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		last_login_at TIMESTAMP NULL,
		INDEX idx_username (username),
		INDEX idx_token (token)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 角色表
	CREATE TABLE IF NOT EXISTS roles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		status TINYINT DEFAULT 1 COMMENT '1:active, 0:inactive',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 权限表
	CREATE TABLE IF NOT EXISTS permissions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		resource VARCHAR(255) NOT NULL COMMENT '资源类型，如: user, role, permission等',
		action VARCHAR(255) NOT NULL COMMENT '操作类型，如: create, read, update, delete',
		status TINYINT DEFAULT 1 COMMENT '1:active, 0:inactive',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_name (name),
		INDEX idx_resource_action (resource, action)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 用户-角色关联表
	CREATE TABLE IF NOT EXISTS user_roles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		role_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by INT COMMENT '创建者ID',
		UNIQUE KEY uk_user_role (user_id, role_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
		INDEX idx_user_id (user_id),
		INDEX idx_role_id (role_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 角色-权限关联表
	CREATE TABLE IF NOT EXISTS role_permissions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		role_id INT NOT NULL,
		permission_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by INT COMMENT '创建者ID',
		UNIQUE KEY uk_role_permission (role_id, permission_id),
		FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
		FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
		INDEX idx_role_id (role_id),
		INDEX idx_permission_id (permission_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 操作日志表
	CREATE TABLE IF NOT EXISTS operation_logs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		action VARCHAR(255) NOT NULL COMMENT '操作类型',
		resource VARCHAR(255) NOT NULL COMMENT '资源类型',
		resource_id INT COMMENT '资源ID',
		detail TEXT COMMENT '详细信息',
		ip_address VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		INDEX idx_user_id (user_id),
		INDEX idx_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`
