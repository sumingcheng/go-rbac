package model

import "time"

// User 用户模型
type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Password  string `gorm:"size:100;not null"`
	Token     string `gorm:"size:255"`
	Roles     []Role `gorm:"many2many:user_roles;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
