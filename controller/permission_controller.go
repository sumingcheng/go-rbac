package controller

import (
	"net/http"
	"rbac/dto"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permissionService *service.PermissionService
}

func NewPermissionController(permissionService *service.PermissionService) *PermissionController {
	return &PermissionController{permissionService: permissionService}
}

func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.permissionService.CreatePermission(req.Name, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Code:    http.StatusCreated,
		Message: "权限创建成功",
	})
}
