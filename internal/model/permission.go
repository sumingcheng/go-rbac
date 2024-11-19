package model

import "time"

// Permission 权限模型
type Permission struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Resource    string    `db:"resource"`
	Action      string    `db:"action"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
