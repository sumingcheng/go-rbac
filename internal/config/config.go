package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Server   struct {
		Port string
	}
}

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func NewConfig() *Config {
	config := &Config{}

	// 从环境变量加载配置，如果没有则使用默认值
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "3306")
	config.Database.Username = getEnv("DB_USER", "root")
	config.Database.Password = getEnv("DB_PASSWORD", "123456")
	config.Database.Name = getEnv("DB_NAME", "rbac_db")

	config.Server.Port = getEnv("SERVER_PORT", "8080")

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
}
