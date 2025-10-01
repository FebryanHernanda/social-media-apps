package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/febryanhernanda/social-media-apps/internal/repositories"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	repo *repositories.UserRepository
	rdb  *redis.Client
}

func NewUserHandler(repo *repositories.UserRepository, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		repo: repo,
		rdb:  rdb,
	}
}

// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.AllUser
// @Failure 500 {object} utils.ErrorResponse
// @Router /user/ [get]
func (h *UserHandler) GetAllUser(ctx *gin.Context) {
	users, err := h.repo.GetAllUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

/* ======================================================================= NOTIFICATIONS */

// @Summary      Get user notifications
// @Description  Retrieve all notifications for the authenticated user
// @ID           get-notifications
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} models.Notifications "Returns list of notifications or empty array"
// @Failure      401 {object} utils.ErrorResponse "Unauthorized"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /user/notifications [get]
func (h *UserHandler) GetNotifications(ctx *gin.Context) {
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

	notif, err := h.repo.GetNotifications(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(notif) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "No notifications yet",
			"data":    []interface{}{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("you have %d notifications", len(notif)),
		"data":    notif,
	})
}

/* ======================================================================= NOTIFICATIONS */

// @Summary      Mark notification as read
// @Description  Mark a specific notification as read for the authenticated user
// @ID           read-notification
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Notification ID"
// @Success 	 200 {object} map[string]interface{} "Notification marked as read successfully"
// @Failure      400 {object} utils.ErrorResponse "Invalid notification ID"
// @Failure      401 {object} utils.ErrorResponse "Unauthorized"
// @Failure      404 {object} utils.ErrorResponse "Notification not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /user/notifications/{id} [patch]
func (h *UserHandler) ReadNotification(ctx *gin.Context) {
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

	notifIDStr := ctx.Param("id")
	notificationID, err := strconv.Atoi(notifIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	updated, err := h.repo.ReadNotification(ctx, notificationID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !updated {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "notification not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notification marked as read",
	})
}

/* ======================================================================= FOLLOWING */

// @Summary      Follow a user
// @Description  Send a follow request to another user
// @ID           follow-user
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID of the user to follow"
// @Success 	 200 {object} map[string]interface{} "Successfully following the user"
// @Failure      400 {object} utils.ErrorResponse "Invalid user ID / Cannot follow yourself / Already following"
// @Failure      401 {object} utils.ErrorResponse "Unauthorized"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /user/{id}/follow [post]
func (h *UserHandler) FollowRequest(ctx *gin.Context) {
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

	followedIDStr := ctx.Param("id")
	followedUserID, err := strconv.Atoi(followedIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	if userID == followedUserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "cannot follow yourself",
		})
	}

	followReq := models.Follows{
		UserID:         userID,
		FollowedUserID: followedUserID,
	}

	follow, err := h.repo.FollowRequest(ctx, &followReq)
	if err != nil {
		if err.Error() == "already following this user" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("You are now following user with ID %d", follow.FollowedUserID),
	})

}

// @Summary      Unfollow a user
// @Description  Remove a follow relationship with another user
// @ID           unfollow-user
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID of the user to unfollow"
// @Success 200 {object} map[string]interface{} "Unfollowed successfully"
// @Failure      400 {object} utils.ErrorResponse "Invalid user ID / Not following this user"
// @Failure      401 {object} utils.ErrorResponse "Unauthorized"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /user/{id}/unfollow [delete]
func (h *UserHandler) UnfollowRequest(ctx *gin.Context) {
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

	followedIDStr := ctx.Param("id")
	followedUserID, err := strconv.Atoi(followedIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	err = h.repo.UnfollowRequest(ctx, userID, followedUserID)
	if err != nil {
		if err.Error() == "not following this user" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "unfollowed successfully",
	})
}
