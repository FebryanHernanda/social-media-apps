package routers

import (
	"github.com/febryanhernanda/social-media-apps/internal/handlers"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(r *gin.Engine, jwtManager *utils.JWTManager, rdb *redis.Client, authHandler *handlers.AuthHandler) {
	authRoutes := r.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)
}
