package router

import (
	"rbac/controller"
	"rbac/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userController *controller.UserController,
	roleController *controller.RoleController,
	permissionController *controller.PermissionController,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	r := gin.Default()

	// 使用 CORS 中间件
	r.Use(middleware.Cors())

	// API 路由分组
	api := r.Group("/api")
	{
		// 公开路由
		api.POST("/login", userController.Login)
		api.POST("/register", userController.CreateUser)

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(authMiddleware.AuthRequired())
		{
			// 需要特定权限的路由
			protected.POST("/roles", authMiddleware.RequirePermission("create_role"), roleController.CreateRole)
			protected.POST("/permissions", authMiddleware.RequirePermission("create_permission"), permissionController.CreatePermission)
			protected.POST("/users/:userID/roles/:roleID", authMiddleware.RequirePermission("assign_role"), userController.AssignRole)
		}
	}

	return r
}
