package routers

import (
	"github.com/febryanhernanda/social-media-apps/internal/handlers"
	"github.com/febryanhernanda/social-media-apps/internal/middlewares"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func PostRouter(r *gin.Engine, postHandler *handlers.PostHandler, jwtManager *utils.JWTManager, rdb *redis.Client) {
	postRoutes := r.Group("/post")
	postRoutes.Use(middlewares.VerifyToken(jwtManager, rdb))
	postRoutes.POST("/", postHandler.CreatePost)
	postRoutes.POST(":id/comment", postHandler.AddComment)
	postRoutes.POST(":id/like", postHandler.LikePost)
	postRoutes.DELETE(":id/unlike", postHandler.UnlikePost)
}
