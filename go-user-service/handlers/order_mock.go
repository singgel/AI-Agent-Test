package handlers

import (
	"go-user-service/models"
	"go-user-service/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MockOrderHandler struct{}

func NewMockOrderHandler() *MockOrderHandler {
	return &MockOrderHandler{}
}

func (h *MockOrderHandler) Create(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	order := &models.Order{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Status:      models.OrderStatusPending,
		Description: req.Description,
	}

	if err := models.MockCreateOrder(order); err != nil {
		utils.InternalServerError(c, "Failed to create order")
		return
	}

	utils.Success(c, order.ToResponse())
}

func (h *MockOrderHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	order, err := models.MockGetOrderByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	utils.Success(c, order.ToResponse())
}

func (h *MockOrderHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) == 0 {
		order, _ := models.MockGetOrderByID(uint(id))
		utils.Success(c, order.ToResponse())
		return
	}

	order, err := models.MockUpdateOrder(uint(id), updates)
	if err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	utils.Success(c, order.ToResponse())
}

func (h *MockOrderHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	if err := models.MockDeleteOrder(uint(id)); err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	utils.Success(c, nil)
}

func (h *MockOrderHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	userIDStr := c.Query("user_id")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	var userID uint
	if userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userID = uint(uid)
		}
	}

	orders, total := models.MockListOrders(page, size, userID, status)

	items := make([]models.OrderResponse, len(orders))
	for i, order := range orders {
		items[i] = order.ToResponse()
	}

	response := models.OrderListResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Items: items,
	}

	utils.Success(c, response)
}
