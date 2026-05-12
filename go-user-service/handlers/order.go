package handlers

import (
	"go-user-service/models"
	"go-user-service/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	order := models.Order{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Status:      models.OrderStatusPending,
		Description: req.Description,
	}

	if err := models.DB.Create(&order).Error; err != nil {
		utils.InternalServerError(c, "Failed to create order")
		return
	}

	order.OrderNo = order.GenerateOrderNo()
	if err := models.DB.Save(&order).Error; err != nil {
		utils.InternalServerError(c, "Failed to generate order number")
		return
	}

	utils.Success(c, order.ToResponse())
}

func (h *OrderHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	utils.Success(c, order.ToResponse())
}

func (h *OrderHandler) Update(c *gin.Context) {
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

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "Order not found")
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

	if len(updates) > 0 {
		if err := models.DB.Model(&order).Updates(updates).Error; err != nil {
			utils.InternalServerError(c, "Failed to update order")
			return
		}
	}

	utils.Success(c, order.ToResponse())
}

func (h *OrderHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	if err := models.DB.Delete(&order).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete order")
		return
	}

	utils.Success(c, nil)
}

func (h *OrderHandler) List(c *gin.Context) {
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

	query := models.DB.Model(&models.Order{})

	if userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.InternalServerError(c, "Failed to count orders")
		return
	}

	var orders []models.Order
	offset := (page - 1) * size
	if err := query.Order("created_at DESC").Limit(size).Offset(offset).Find(&orders).Error; err != nil {
		utils.InternalServerError(c, "Failed to query orders")
		return
	}

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
