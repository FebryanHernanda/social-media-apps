package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/febryanhernanda/social-media-apps/internal/repositories"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type FeedHandler struct {
	repo *repositories.FeedRepository
	rdb  *redis.Client
}

func NewFeedHandler(repo *repositories.FeedRepository, rdb *redis.Client) *FeedHandler {
	return &FeedHandler{
		repo: repo,
		rdb:  rdb,
	}
}

// @Summary Get user feed
// @Description Get feed of a user
// @ID get-user-feed
// @Tags feed
// @Security     BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.FeedPost
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /feed [get]
func (h *FeedHandler) GetUserFeed(ctx *gin.Context) {
	rawClaims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}
	claims := rawClaims.(*utils.Claims)
	userID := claims.UserID

	redisKey := "feed:post"
	var cached []models.FeedPost
	if h.rdb != nil {
		err := utils.GetCache(ctx, h.rdb, redisKey, &cached)
		if err != nil {
			log.Panicln("Redis error, back to DB : ", err)
		}
		if len(cached) > 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    cached,
				"message": "data from cache",
			})
			return
		}
	}

	feed, err := h.repo.GetUserFeed(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(feed) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "No posts found. Follow some users to see their posts.",
			"data":    []interface{}{},
		})
		return
	}

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, feed, 60*time.Second)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    feed,
	})
}
