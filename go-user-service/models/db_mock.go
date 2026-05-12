package models

import (
	"errors"
	"sync"
	"time"
)

var (
	mockUsers     = make(map[uint]*User)
	mockUsersByUsername = make(map[string]*User)
	mockUsersByEmail    = make(map[string]*User)
	mockOrders    = make(map[uint]*Order)
	mockNextID    uint = 1
	mockOrderNextID uint = 1
	mockMutex     sync.RWMutex
)

type MockDB struct{}

func InitMockDB() error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	mockUsers = make(map[uint]*User)
	mockUsersByUsername = make(map[string]*User)
	mockUsersByEmail = make(map[string]*User)
	mockOrders = make(map[uint]*Order)
	mockNextID = 1
	mockOrderNextID = 1

	return nil
}

func MockCreateUser(user *User) error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	user.ID = mockNextID
	mockNextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	mockUsers[user.ID] = user
	mockUsersByUsername[user.Username] = user
	mockUsersByEmail[user.Email] = user

	return nil
}

func MockGetUserByID(id uint) (*User, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	user, exists := mockUsers[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	userCopy := *user
	return &userCopy, nil
}

func MockGetUserByUsername(username string) (*User, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	user, exists := mockUsersByUsername[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	userCopy := *user
	return &userCopy, nil
}

func MockGetUserByEmail(email string) (*User, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	user, exists := mockUsersByEmail[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	userCopy := *user
	return &userCopy, nil
}

func MockCheckUsernameExists(username string) bool {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	_, exists := mockUsersByUsername[username]
	return exists
}

func MockCheckEmailExists(email string) bool {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	_, exists := mockUsersByEmail[email]
	return exists
}

// Mock Order Functions

func MockCreateOrder(order *Order) error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	order.ID = mockOrderNextID
	mockOrderNextID++
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.OrderNo = time.Now().Format("20060102150405") + string(rune(order.ID%1000+48))

	mockOrders[order.ID] = order
	return nil
}

func MockGetOrderByID(id uint) (*Order, error) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	order, exists := mockOrders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	orderCopy := *order
	return &orderCopy, nil
}

func MockUpdateOrder(id uint, updates map[string]interface{}) (*Order, error) {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	order, exists := mockOrders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	if amount, ok := updates["amount"].(float64); ok {
		order.Amount = amount
	}
	if status, ok := updates["status"].(OrderStatus); ok {
		order.Status = status
	}
	if description, ok := updates["description"].(string); ok {
		order.Description = description
	}
	order.UpdatedAt = time.Now()

	orderCopy := *order
	return &orderCopy, nil
}

func MockDeleteOrder(id uint) error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	if _, exists := mockOrders[id]; !exists {
		return errors.New("order not found")
	}

	delete(mockOrders, id)
	return nil
}

func MockListOrders(page, size int, userID uint, status string) ([]Order, int64) {
	mockMutex.RLock()
	defer mockMutex.RUnlock()

	var orders []Order
	for _, order := range mockOrders {
		if userID > 0 && order.UserID != userID {
			continue
		}
		if status != "" && string(order.Status) != status {
			continue
		}
		orders = append(orders, *order)
	}

	total := int64(len(orders))

	// Simple pagination
	start := (page - 1) * size
	if start >= len(orders) {
		return []Order{}, total
	}
	end := start + size
	if end > len(orders) {
		end = len(orders)
	}

	return orders[start:end], total
}
