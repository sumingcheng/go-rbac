package model

import "time"

// User 用户模型
type User struct {
    ID        string    `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:100;uniqueIndex;not null"`
    Password  string    `gorm:"size:100;not null"`
    Roles     []Role    `gorm:"many2many:user_roles;"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Role 角色模型
type Role struct {
    ID          string       `gorm:"primaryKey"`
    Name        string       `gorm:"size:100;uniqueIndex;not null"`
    Description string       `gorm:"size:200"`
    Permissions []Permission `gorm:"many2many:role_permissions;"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Permission 权限模型
type Permission struct {
    ID          string    `gorm:"primaryKey"`
    Name        string    `gorm:"size:100;uniqueIndex;not null"`
    Description string    `gorm:"size:200"`
    Resource    string    `gorm:"size:100;not null"` // 资源，如：user, role, article
    Action      string    `gorm:"size:100;not null"` // 操作，如：create, read, update, delete
    CreatedAt   time.Time
    UpdatedAt   time.Time
} 