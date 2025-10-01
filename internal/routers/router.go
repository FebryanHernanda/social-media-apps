package routers

import (
	"log"
	"os"

	"github.com/febryanhernanda/social-media-apps/docs"
	"github.com/febryanhernanda/social-media-apps/internal/handlers"
	"github.com/febryanhernanda/social-media-apps/internal/repositories"
	"github.com/febryanhernanda/social-media-apps/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	/* JWT */
	jwtSecret := os.Getenv("JWTKEY")
	if jwtSecret == "" {
		log.Fatal("JWT Key env variable not set")
	}
	jwtManager := utils.NewJWTManager(jwtSecret)

	/* Repo & Handler */
	authRepo := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepo, jwtManager, rdb)

	postRepo := repositories.NewPostRepository(db)
	postHandler := handlers.NewPostHandler(postRepo, rdb)

	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo, rdb)

	feedRepo := repositories.NewFeedRepository(db)
	feedHandler := handlers.NewFeedHandler(feedRepo, rdb)

	/* Register Router */
	AuthRouter(r, jwtManager, rdb, authHandler)
	PostRouter(r, postHandler, jwtManager, rdb)
	UserRouter(r, userHandler, jwtManager, rdb)
	FeedRouter(r, feedHandler, jwtManager, rdb)

	/* Register file upload */
	r.Static("/public", "./public")

	/* Register Swagger */
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
