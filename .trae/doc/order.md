帮我创建订单服务模块

# 技术要求
- 框架：Gin + GORM
- 目录：internal/

# 文件结构
- internal/model/order.go - 订单模型
- internal/service/order.go - 业务逻辑
- internal/handler/order.go - HTTP接口
- internal/repository/order.go - 数据访问

# 接口规格
1. POST /api/v1/orders - 创建订单
   - 请求：{ user_id, product_id, quantity }
   - 响应：{ order_id, status, created_at }

2. GET /api/v1/orders/:id - 查询订单
   - 响应：{ order详情 }

3. GET /api/v1/orders - 订单列表
   - 支持分页：page, page_size
   - 支持筛选：status, user_id

4. PUT /api/v1/orders/:id/status - 更新状态
   - 请求：{ status }

# 约束
- 参考 #internal/handler/user.go 的代码风格
- 使用统一的响应格式
- 错误码定义在 pkg/errors/ 下

