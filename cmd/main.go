package main

import (
	"log"

	"github.com/febryanhernanda/social-media-apps/internal/configs"
	"github.com/febryanhernanda/social-media-apps/internal/routers"
)

// @title Social Media API
// @version 1.0
// @description API for Social Media
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db, err := configs.InitDB()
	if err != nil {
		log.Fatal("DB init failed: ", err)
	}
	defer db.Close()

	rdb, err := configs.InitRedis()
	if err != nil {
		log.Fatal("RDB init failed: ", err)
	}

	if rdb != nil {
		defer rdb.Close()
	}

	r := routers.Router(db, rdb)

	r.Run(":8080")
}
