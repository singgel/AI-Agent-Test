package models

import (
	"time"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderNo     string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	Amount      float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status      OrderStatus    `json:"status" gorm:"size:20;default:'pending'"`
	Description string         `json:"description" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateOrderRequest struct {
	UserID      uint    `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"max=500"`
}

type UpdateOrderRequest struct {
	Amount      float64     `json:"amount" binding:"omitempty,gt=0"`
	Status      OrderStatus `json:"status" binding:"omitempty,oneof=pending paid shipped completed cancelled"`
	Description string      `json:"description" binding:"max=500"`
}

type OrderResponse struct {
	ID          uint        `json:"id"`
	OrderNo     string      `json:"order_no"`
	UserID      uint        `json:"user_id"`
	Amount      float64     `json:"amount"`
	Status      OrderStatus `json:"status"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderListResponse struct {
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Items []OrderResponse `json:"items"`
}

func (o *Order) ToResponse() OrderResponse {
	return OrderResponse{
		ID:          o.ID,
		OrderNo:     o.OrderNo,
		UserID:      o.UserID,
		Amount:      o.Amount,
		Status:      o.Status,
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func (o *Order) GenerateOrderNo() string {
	return time.Now().Format("20060102150405") + string(rune(o.ID%1000))
}
