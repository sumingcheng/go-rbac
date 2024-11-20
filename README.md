

# Go RBAC 系统

基于 Go 语言实现的角色权限管理系统 (Role-Based Access Control)。

## 项目架构

```
├── main.go # 应用程序入口
├── controller/ # 控制器层
│ ├── user.go # 用户相关控制器
│ ├── role.go # 角色相关控制器
│ └── permission.go # 权限相关控制器
├── service/ # 业务逻辑层
│ ├── user.go # 用户服务
│ ├── role.go # 角色服务
│ └── permission.go # 权限服务
├── model/ # 数据模型
│ ├── user.go # 用户模型
│ ├── role.go # 角色模型
│ └── permission.go # 权限模型
├── dto/ # 数据传输对象
│ └── response.go # 响应结构定义
├── internal/ # 内部包
│ ├── config/ # 配置管理
│ └── config.go # 配置结构和方法
│ ├── database/ # 数据库相关
│ │ ├── database.go # 数据库初始化
│ │ ├── schema.go # 数据库表结构
│ │ └── seed.go # 初始数据填充
├── middleware/ # 中间件
│ ├── auth.go # 认证中间件
│ └── cors.go # CORS 中间件
└── router/ # 路由配置
└── router.go # 路由注册
```

## 核心功能

### 1. 用户认证
- 用户注册与登录
- 基于 Token 的认证机制
- 密码加密存储

### 2. 权限控制
- 基于 RBAC 模型的权限管理
- 角色-权限分配
- 用户-角色分配
- 权限验证中间件

## 技术栈

- Web 框架：Gin
- 数据库：MySQL
- ORM：sqlx
- 认证：Token-based Authentication

## 快速开始

1. 克隆项目
2. 配置环境变量

项目地址：[GitHub - sumingcheng/go-rbac](https://github.com/sumingcheng/go-rbac)
