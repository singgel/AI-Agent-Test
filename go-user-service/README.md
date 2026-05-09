# Go User Service

使用 Gin + GORM + MySQL 构建的用户服务 HTTP API。

## 功能特性

- 用户注册
- 用户登录
- 获取用户信息
- JWT 身份认证
- 密码加密存储

## 项目结构

```
go-user-service/
├── config/          # 配置文件
│   ├── config.yaml  # 配置文件
│   └── config.go    # 配置加载
├── handlers/        # HTTP 处理器
│   └── user.go      # 用户相关接口
├── models/          # 数据模型
│   ├── db.go        # 数据库连接
│   └── user.go      # 用户模型
├── routes/          # 路由配置
│   └── routes.go    # 路由设置
├── utils/           # 工具函数
│   ├── jwt.go       # JWT 工具
│   ├── password.go  # 密码加密
│   └── response.go  # 响应封装
├── main.go          # 入口文件
├── go.mod           # Go 模块
├── start.sh         # 启动脚本
└── README.md        # 项目说明
```

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 5.7+

### 安装依赖

```bash
go mod download
```

### 配置数据库

编辑 `config/config.yaml` 文件，修改数据库连接信息：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  dbname: user_service
```

### 启动服务

使用启动脚本：

```bash
chmod +x start.sh
./start.sh
```

或直接运行：

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 1. 用户注册

**POST** `/api/user/register`

请求体：
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 用户登录

**POST** `/api/user/login`

请求体：
```json
{
  "username": "testuser",
  "password": "password123"
}
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 3. 获取用户信息

**GET** `/api/user/:id`

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 1001 | 用户名已存在 |
| 1002 | 邮箱已存在 |
| 1003 | 用户名或密码错误 |

## 测试示例

### 注册
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### 登录
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 获取用户信息
```bash
curl -X GET http://localhost:8080/api/user/1
```
