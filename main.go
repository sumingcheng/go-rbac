package main

import (
	"rbac/controller"
	"rbac/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "root:admin123456@tcp(172.22.220.64:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func main() {
	db := initDB() // 初始化数据库连接

	rbacService := service.NewRBACService(db)
	rbacController := controller.NewRBACController(rbacService)

	r := gin.Default()

	// 用户相关路由
	r.POST("/users", rbacController.CreateUser)
	r.POST("/users/:userID/roles/:roleID", rbacController.AssignRole)

	r.Run(":8080")
}
