package routers

import (
	"github.com/febryanhernanda/social-media-apps/internal/handlers"
	"github.com/febryanhernanda/social-media-apps/internal/middlewares"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func FeedRouter(r *gin.Engine, feedHandler *handlers.FeedHandler, jwtManager *utils.JWTManager, rdb *redis.Client) {
	feedRoutes := r.Group("/feed")
	feedRoutes.Use(middlewares.VerifyToken(jwtManager, rdb))
	feedRoutes.GET("/", feedHandler.GetUserFeed)
}
