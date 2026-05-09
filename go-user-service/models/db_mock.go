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
	mockNextID    uint = 1
	mockMutex     sync.RWMutex
)

type MockDB struct{}

func InitMockDB() error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	mockUsers = make(map[uint]*User)
	mockUsersByUsername = make(map[string]*User)
	mockUsersByEmail = make(map[string]*User)
	mockNextID = 1

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
