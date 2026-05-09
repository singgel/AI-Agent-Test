package handlers

import (
	"go-user-service/config"
	"go-user-service/models"
	"go-user-service/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MockUserHandler struct {
	jwtUtil *utils.JWTUtil
}

func NewMockUserHandler(cfg *config.JWTConfig) *MockUserHandler {
	models.InitMockDB()
	return &MockUserHandler{
		jwtUtil: utils.NewJWTUtil(cfg.Secret, cfg.ExpireHours),
	}
}

func (h *MockUserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if models.MockCheckUsernameExists(req.Username) {
		utils.Error(c, 1001, "Username already exists")
		return
	}

	if models.MockCheckEmailExists(req.Email) {
		utils.Error(c, 1002, "Email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password")
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := models.MockCreateUser(user); err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	utils.Success(c, user.ToResponse())
}

func (h *MockUserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := models.MockGetUserByUsername(req.Username)
	if err != nil {
		utils.Error(c, 1003, "Invalid username or password")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.Error(c, 1003, "Invalid username or password")
		return
	}

	token, err := h.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	utils.Success(c, response)
}

func (h *MockUserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := models.MockGetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user.ToResponse())
}

func (h *MockUserHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := h.jwtUtil.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func (h *MockUserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := models.MockGetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user.ToResponse())
}
