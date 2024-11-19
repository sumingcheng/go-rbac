package main

import (
	"rbac/controller"
	"rbac/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "root:your_password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
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
