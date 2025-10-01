package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/febryanhernanda/social-media-apps/internal/repositories"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PostHandler struct {
	repo *repositories.PostRepository
	rdb  *redis.Client
}

func NewPostHandler(repo *repositories.PostRepository, rdb *redis.Client) *PostHandler {
	return &PostHandler{
		repo: repo,
		rdb:  rdb,
	}
}

// @Summary      Create a post
// @Description  Create a post with optional image and text content
// @ID           create-post
// @Tags         post
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        content formData string true "Post content"
// @Param        image   formData file false "Post image file"
// @Success      200 {object} models.Post "Post created successfully"
// @Failure      400 {object} utils.ErrorResponse "Bad request"
// @Failure      401 {object} utils.ErrorResponse "Unauthorized"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /post [post]
func (h *PostHandler) CreatePost(ctx *gin.Context) {
	rawClaims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}
	claims := rawClaims.(*utils.Claims)

	var req models.CreatePostRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var filePath *string = nil

	_, err := ctx.FormFile("image")
	if err == nil {
		uploadedPath, err := utils.UploadFile(ctx, "image", "public/post", "post", "post")
		if err != nil {
			log.Printf("[DEBUG] ERRORS : %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to upload photo",
			})
			return
		}
		filePath = &uploadedPath
	}
	post := &models.Post{
		Content:   req.Content,
		ImagePath: filePath,
		UserID:    claims.UserID,
	}

	newPost, err := h.repo.CreatePost(ctx, post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create post",
		})
		return
	}

	if err := utils.InvalidateCache(ctx, h.rdb, []string{"feed:post"}); err != nil {
		log.Println("Redis delete cache error:", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "post created successfully",
		"data":    newPost,
	})
}

/* ======================================================================= LIKE POST */

// @Summary Like a post
// @Description Like a post with user ID
// @ID like-post
// @Tags post
// @Security     BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "post ID"
// @Success 200 {object} models.Like
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /post/{id}/like [post]
func (h *PostHandler) LikePost(ctx *gin.Context) {
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

	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	postOwnerID, err := h.repo.GetPostOwnerID(ctx, postID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := utils.InvalidateCache(ctx, h.rdb, []string{"feed:post"}); err != nil {
		log.Println("Redis delete cache error:", err)
	}

	like, err := h.repo.LikePost(ctx, postID, userID, postOwnerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("post id %s liked", postIDStr),
		"data":    like,
	})
}

// @Summary Unlike a post
// @Description Unlike a post with user ID
// @ID unlike-post
// @Tags post
// @Security     BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "post ID"
// @Success 200 {object} models.Like
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /post/{id}/unlike [delete]
func (h *PostHandler) UnlikePost(ctx *gin.Context) {
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

	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	err = h.repo.UnlikePost(ctx, postID, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := utils.InvalidateCache(ctx, h.rdb, []string{"feed:post"}); err != nil {
		log.Println("Redis delete cache error:", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("post id %s unliked", postIDStr),
	})
}

/* ======================================================================= COMMENT POST */

// @Summary Add a comment to a post
// @Description Add a comment to a post with the given content
// @ID add-comment
// @Tags post
// @Security     BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "post ID"
// @Param body body models.CommentRequest true "comment request body"
// @Success 200 {object} models.Comment
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /post/{id}/comment [post]
func (h *PostHandler) AddComment(ctx *gin.Context) {
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

	postIDStr := ctx.Param("id")
	postID, _ := strconv.Atoi(postIDStr)

	var req models.CommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	postOwnerID, err := h.repo.GetPostOwnerID(ctx, postID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "post not found",
		})
		return
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	comment, err = h.repo.AddComment(ctx, comment, postOwnerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := utils.InvalidateCache(ctx, h.rdb, []string{"feed:post"}); err != nil {
		log.Println("Redis delete cache error:", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "comment added",
		"data":    comment,
	})
}
