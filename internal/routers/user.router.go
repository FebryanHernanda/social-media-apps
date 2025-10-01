package routers

import (
	"github.com/febryanhernanda/social-media-apps/internal/handlers"
	"github.com/febryanhernanda/social-media-apps/internal/middlewares"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func UserRouter(r *gin.Engine, userHandler *handlers.UserHandler, jwtManager *utils.JWTManager, rdb *redis.Client) {
	userRoutes := r.Group("/user")
	userRoutes.GET("/", userHandler.GetAllUser)
	userRoutes.Use(middlewares.VerifyToken(jwtManager, rdb))
	userRoutes.POST("/:id/follow", userHandler.FollowRequest)
	userRoutes.DELETE("/:id/unfollow", userHandler.UnfollowRequest)

	userRoutes.GET("/notifications", userHandler.GetNotifications)
	userRoutes.PATCH("/notifications/:id", userHandler.ReadNotification)
}
