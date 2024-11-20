package service

import (
	"rbac/internal/middleware"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sqlx.DB
}

// 确保 UserService 实现了 AuthService 接口
var _ middleware.AuthService = (*UserService)(nil)

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		username, string(hashedPassword),
	)
	return err
}

func (s *UserService) AssignRole(userID, roleID int) error {
	_, err := s.db.Exec(
		"INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)",
		userID, roleID,
	)
	return err
}

func (s *UserService) HasPermission(userID int, permissionName string) (bool, error) {
	var count int
	err := s.db.Get(&count, `
        SELECT COUNT(*) FROM users u 
        JOIN user_roles ur ON u.id = ur.user_id 
        JOIN roles r ON ur.role_id = r.id 
        JOIN role_permissions rp ON r.id = rp.role_id 
        JOIN permissions p ON rp.permission_id = p.id 
        WHERE u.id = ? AND p.name = ?`,
		userID, permissionName,
	)
	return count > 0, err
}

func (s *UserService) ValidateToken(token string) (int, error) {
	var userID int
	err := s.db.Get(&userID, "SELECT id FROM users WHERE token = ?", token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *UserService) CheckPermission(userID int, permission string) (bool, error) {
	var count int
	err := s.db.Get(&count, `
		SELECT COUNT(*) FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE u.id = ? AND p.name = ?
	`, userID, permission)
	return count > 0, err
}

func (s *UserService) Login(username, password string) (string, error) {
	var user struct {
		ID       int
		Password string
	}

	err := s.db.Get(&user, "SELECT id, password FROM users WHERE username = ?", username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	// 生成token
	token := generateToken() // 你需要实现这个函数，可以使用UUID或其他方式

	// 更新用户token
	_, err = s.db.Exec("UPDATE users SET token = ? WHERE id = ?", token, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken() string {
	// 这里可以使用UUID或其他方式生成token
	return uuid.New().String()
}
