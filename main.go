package main

import (
	"rbac/controller"
	"rbac/internal/config"
	"rbac/internal/database"
	"rbac/internal/router"
	"rbac/service"
)

func main() {
	cfg := config.NewConfig()
	db := database.InitDB(cfg)

	// 初始化 services
	userService := service.NewUserService(db)
	roleService := service.NewRoleService(db)
	permissionService := service.NewPermissionService(db)

	// 初始化 controllers
	userController := controller.NewUserController(userService)
	roleController := controller.NewRoleController(roleService)
	permissionController := controller.NewPermissionController(permissionService)

	// 设置路由
	r := router.SetupRouter(
		userController,
		roleController,
		permissionController,
	)

	r.Run(":" + cfg.Server.Port)
}
