package handlers

import (
	"log"
	"net/http"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/febryanhernanda/social-media-apps/internal/repositories"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo       *repositories.AuthRepository
	JWTManager *utils.JWTManager
	rdb        *redis.Client
}

func NewAuthHandler(repo *repositories.AuthRepository, jwtManager *utils.JWTManager, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{
		repo:       repo,
		JWTManager: jwtManager,
		rdb:        rdb,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req models.RegisterUser
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to hash password",
		})
		return
	}
	req.Password = string(hashedPass)

	user, err := h.repo.RegisterUser(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "registration has failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "registration success",
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var user models.LoginUser

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	userFromDB, err := h.repo.LoginUser(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "invalid email or password",
		})
		return
	}

	token, err := h.JWTManager.GenerateToken(userFromDB)
	if err != nil {
		log.Printf("[DEBUG] Error : %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successfully",
		"token":   token,
	})
}
