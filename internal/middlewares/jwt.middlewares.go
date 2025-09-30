package middlewares

import (
	"net/http"
	"strings"

	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

/* Verify token from middlewares */
func VerifyToken(jwtManager *utils.JWTManager, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		parts := strings.Fields(authHeader)

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Please login",
			})
			return
		}

		tokenString := parts[1]

		/* Check Blacklist token on redis */
		exists, err := rdb.Exists(ctx, "blacklist:"+tokenString).Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Server error",
			})
			return
		}

		if exists > 0 {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Token Invalid, Please login again",
			})
			return
		}

		/* Check and Validate expire token */
		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
