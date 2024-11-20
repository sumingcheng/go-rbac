package database

import (
	"fmt"
	"log"
	"rbac/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB(cfg *config.Config) *sqlx.DB {
	// 先创建数据库连接（不指定数据库名）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 创建数据库（如果不存在）
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.Database.Name)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	// 关闭当前连接
	db.Close()

	// 使用新的数据库重新连接
	db, err = sqlx.Connect("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 初始化表结构
	if err := InitSchema(db); err != nil {
		log.Fatalf("初始化数据库表结构失败: %v", err)
	}
	return db
}
