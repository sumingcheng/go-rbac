package database

import (
	"fmt"
	"log"
	"rbac/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB(cfg *config.Config) *sqlx.DB {
	// 先创建数据库连接（不指定数据库名）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 创建数据库（如果不存在）
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.Database.Name)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	// 关闭当前连接
	db.Close()

	// 使用新的数据库重新连接
	db, err = sqlx.Connect("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 初始化表结构
	initSchema(db)
	return db
}

func initSchema(db *sqlx.DB) {
	schema := `
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

	// 执行建表语句
	db.MustExec(schema)

	// 初始化基础数据
	initBaseData(db)
}

// 初始化基础数据
func initBaseData(db *sqlx.DB) {
	// 创建超级管理员角色
	db.MustExec(`
        INSERT IGNORE INTO roles (name, description) 
        VALUES ('super_admin', '超级管理员');
    `)

	// 创建基础权限
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
		db.MustExec(`
            INSERT IGNORE INTO permissions (name, description, resource, action) 
            VALUES (?, ?, ?, ?)
        `, perm.Name, perm.Description, perm.Resource, perm.Action)
	}

	// 为超级管理员分配所有权限
	db.MustExec(`
        INSERT IGNORE INTO role_permissions (role_id, permission_id)
        SELECT r.id, p.id
        FROM roles r
        CROSS JOIN permissions p
        WHERE r.name = 'super_admin'
    `)
}
