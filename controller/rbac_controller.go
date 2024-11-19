package controller

import (
	"rbac/dto"
	"rbac/model"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

type RBACController struct {
	rbacService *service.RBACService
}

func NewRBACController(rbacService *service.RBACService) *RBACController {
	return &RBACController{
		rbacService: rbacService,
	}
}

// CreateUser 创建用户接口
func (c *RBACController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, dto.Response{
			Code:    400,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	userDTO, err := c.rbacService.CreateUser(&user)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Code:    200,
		Message: "Success",
		Data:    userDTO,
	})
}

// AssignRole 分配角色接口
func (c *RBACController) AssignRole(ctx *gin.Context) {
	userID := ctx.Param("userID")
	roleID := ctx.Param("roleID")

	err := c.rbacService.AssignRoleToUser(userID, roleID)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Code:    200,
		Message: "Role assigned successfully",
	})
}
