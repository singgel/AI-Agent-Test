#!/bin/bash

# 用户服务 Mock 模式启动脚本（无需数据库）

# 设置环境变量
export GO_ENV=development

echo "=========================================="
echo "Starting server in MOCK mode"
echo "No database required - using in-memory storage"
echo "=========================================="
echo ""

# 下载依赖
echo "Downloading dependencies..."
go mod download

# 运行服务（使用 mock 版本）
echo "Starting server..."
go run main.go
