package routes

import (
	"go-user-service/config"
	"go-user-service/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	userHandler := handlers.NewUserHandler(&cfg.JWT)

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.GET("/:id", userHandler.GetUser)
		}
	}
}
