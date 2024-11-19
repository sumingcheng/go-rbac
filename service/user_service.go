package service

import (
	"rbac/dto"
	"rbac/model"

	"gorm.io/gorm"
)

type RBACService struct {
	db *gorm.DB
}

func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{db: db}
}

// CreateUser 创建用户
func (s *RBACService) CreateUser(user *model.User) (*dto.UserDTO, error) {
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return s.convertToUserDTO(user), nil
}

// AssignRoleToUser 为用户分配角色
func (s *RBACService) AssignRoleToUser(userID string, roleID string) error {
	var user model.User
	var role model.Role

	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	return s.db.Model(&user).Association("Roles").Append(&role)
}

// CheckPermission 检查用户是否有某个权限
func (s *RBACService) CheckPermission(userID string, resource string, action string) bool {
	var user model.User
	if err := s.db.Preload("Roles.Permissions").First(&user, "id = ?", userID).Error; err != nil {
		return false
	}

	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Resource == resource && perm.Action == action {
				return true
			}
		}
	}
	return false
}

// 转换方法
func (s *RBACService) convertToUserDTO(user *model.User) *dto.UserDTO {
	userDTO := &dto.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	for _, role := range user.Roles {
		roleDTO := dto.RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
		}
		userDTO.Roles = append(userDTO.Roles, roleDTO)
	}

	return userDTO
}
