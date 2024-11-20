package controller

import (
	"net/http"
	"rbac/dto"
	"rbac/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService *service.RoleService
}

func NewRoleController(roleService *service.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.roleService.CreateRole(req.Name, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Code:    http.StatusCreated,
		Message: "角色创建成功",
	})
}

func (c *RoleController) AssignPermission(ctx *gin.Context) {
	roleID, _ := strconv.Atoi(ctx.Param("roleID"))
	permissionID, _ := strconv.Atoi(ctx.Param("permissionID"))

	if err := c.roleService.AssignPermission(roleID, permissionID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "权限分配成功",
	})
}
