package controller

import (
	"rbac/dto"
	"rbac/model"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser 创建用户接口
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, dto.Response{
			Code:    400,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	userDTO, err := c.userService.CreateUser(&user)
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
func (c *UserController) AssignRole(ctx *gin.Context) {
	userID := ctx.Param("userID")
	roleID := ctx.Param("roleID")

	err := c.userService.AssignRoleToUser(userID, roleID)
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
