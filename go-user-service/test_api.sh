#!/bin/bash

# API 测试脚本

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "Testing User Service API"
echo "=========================================="
echo ""

# 测试 1: 注册用户
echo "1. Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/user/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}')
echo "Response: $REGISTER_RESPONSE"
echo ""

# 测试 2: 重复注册（应该失败）
echo "2. Testing Duplicate Registration (should fail)..."
DUPLICATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/user/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}')
echo "Response: $DUPLICATE_RESPONSE"
echo ""

# 测试 3: 用户登录
echo "3. Testing User Login..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}')
echo "Response: $LOGIN_RESPONSE"
echo ""

# 提取 Token（简单提取，实际使用 JSON 解析器更好）
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Extracted Token: $TOKEN"
echo ""

# 测试 4: 获取用户信息
echo "4. Testing Get User Info..."
USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/user/1")
echo "Response: $USER_RESPONSE"
echo ""

# 测试 5: 获取不存在的用户
echo "5. Testing Get Non-existent User (should return 404)..."
NOTFOUND_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/user/999")
echo "Response: $NOTFOUND_RESPONSE"
echo ""

# 测试 6: 错误密码登录
echo "6. Testing Login with Wrong Password (should fail)..."
WRONG_PASSWORD_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"wrongpassword"}')
echo "Response: $WRONG_PASSWORD_RESPONSE"
echo ""

# 测试 7: 注册第二个用户
echo "7. Testing Register Second User..."
REGISTER_RESPONSE2=$(curl -s -X POST "${BASE_URL}/api/user/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser2","email":"test2@example.com","password":"password456"}')
echo "Response: $REGISTER_RESPONSE2"
echo ""

# 测试 8: 获取第二个用户信息
echo "8. Testing Get Second User Info..."
USER_RESPONSE2=$(curl -s -X GET "${BASE_URL}/api/user/2")
echo "Response: $USER_RESPONSE2"
echo ""

echo "=========================================="
echo "All tests completed!"
echo "=========================================="
