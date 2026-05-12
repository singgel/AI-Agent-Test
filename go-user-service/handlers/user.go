package handlers

import (
	"go-user-service/config"
	"go-user-service/models"
	"go-user-service/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserHandler struct {
	jwtUtil *utils.JWTUtil
}

func NewUserHandler(cfg *config.JWTConfig) *UserHandler {
	return &UserHandler{
		jwtUtil: utils.NewJWTUtil(cfg.Secret, cfg.ExpireHours),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	var existingUser models.User
	if err := models.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, 1001, "Username already exists")
		return
	}

	if err := models.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, 1002, "Email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	utils.Success(c, user.ToResponse())
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 记录登录请求
	utils.Logger.Info("用户尝试登录", "username", req.Username)

	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// 记录登录失败 - 用户名不存在
		utils.Logger.Warn("登录失败：用户名不存在", "username", req.Username)
		utils.Error(c, 1003, "Invalid username or password")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		// 记录登录失败 - 密码错误
		utils.Logger.Warn("登录失败：密码错误", "username", req.Username, "userID", user.ID)
		utils.Error(c, 1003, "Invalid username or password")
		return
	}

	token, err := h.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		// 记录Token生成失败
		utils.Logger.Error("Token生成失败", "username", req.Username, "userID", user.ID, "error", err.Error())
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	// 记录登录成功
	utils.Logger.Info("用户登录成功", "username", user.Username, "userID", user.ID)

	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	utils.Success(c, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user.ToResponse())
}

func (h *UserHandler) AuthMiddleware() gin.HandlerFunc {
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

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user.ToResponse())
}
