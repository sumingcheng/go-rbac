# 设置项目基本变量
APP_NAME=go-rbac
DOCKER_IMAGE=$(APP_NAME)
MAIN_FILE=main.go
VERSION=v1.0.0

# Go 相关命令和参数
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=gofmt

# Docker 相关命令
DOCKER=docker

# 默认目标
.DEFAULT_GOAL := help

.PHONY: build run test clean docker-build docker-run fmt vet tidy help tag

# 构建应用
build:
	$(GOBUILD) -o $(APP_NAME) $(MAIN_FILE)

# 运行应用
run:
	$(GORUN) $(MAIN_FILE)

# 运行测试
test:
	$(GOTEST) -v ./...

# 代码格式化
fmt:
	$(GOFMT) -w .

# 代码检查
vet:
	$(GOVET) ./...

# 更新依赖
tidy:
	$(GOMOD) tidy

# 清理构建文件
clean:
	rm -f $(APP_NAME)
	rm -rf dist/

# Docker 相关命令
docker-build:
	$(DOCKER) build -t $(DOCKER_IMAGE) -f deploy/Dockerfile .

docker-run:
	$(DOCKER) run -p 8080:8080 $(DOCKER_IMAGE)

# 添加版本标签命令
tag:
	git tag $(VERSION)
	git push origin $(VERSION)

# 帮助信息
help:
	@echo "可用的 make 命令："
	@echo "make build        - 构建应用 (包含版本信息)"
	@echo "make run         - 运行应用"
	@echo "make test        - 运行测试"
	@echo "make fmt         - 格式化代码"
	@echo "make vet         - 代码静态检查"
	@echo "make tidy        - 整理并更新依赖"
	@echo "make clean       - 清理构建文件"
	@echo "make docker-build- 构建 Docker 镜像"
	@echo "make docker-run  - 运行 Docker 容器"
	@echo "make tag         - 创建版本标签"
