package routes

import (
	"go-user-service/config"
	"go-user-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupMockRoutes(r *gin.Engine, cfg *config.Config) {
	userHandler := handlers.NewMockUserHandler(&cfg.JWT)
	orderHandler := handlers.NewMockOrderHandler()

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.GET("/:id", userHandler.GetUser)
		}

		order := api.Group("/order")
		{
			order.POST("", orderHandler.Create)
			order.GET("/:id", orderHandler.Get)
			order.PUT("/:id", orderHandler.Update)
			order.DELETE("/:id", orderHandler.Delete)
			order.GET("", orderHandler.List)
		}
	}
}
