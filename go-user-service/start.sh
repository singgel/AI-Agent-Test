#!/bin/bash

# 用户服务启动脚本

# 设置环境变量
export GO_ENV=development

# 检查 MySQL 是否运行
echo "Checking MySQL connection..."
if ! mysql -h localhost -u root -proot -e "SELECT 1" > /dev/null 2>&1; then
    echo "Warning: Cannot connect to MySQL. Please ensure MySQL is running."
    echo "You may need to:"
    echo "  1. Start MySQL service"
    echo "  2. Update config/config.yaml with correct database credentials"
    echo ""
fi

# 检查数据库是否存在，如果不存在则创建
echo "Checking database..."
mysql -h localhost -u root -proot -e "CREATE DATABASE IF NOT EXISTS user_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || true

# 下载依赖
echo "Downloading dependencies..."
go mod download

# 运行服务
echo "Starting server..."
go run main.go
