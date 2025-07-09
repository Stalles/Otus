package main

import (
	"log"
	"os"
	"socialNetworkOtus/internal/api"
	"socialNetworkOtus/internal/config"
	"socialNetworkOtus/internal/db"
	"socialNetworkOtus/internal/handlers"
	"socialNetworkOtus/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	connStr := "user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " host=" + cfg.DBHost + " port=" + cfg.DBPort + " sslmode=disable"
	if err := db.RunMigrations(connStr, "migrations/001_create_users.sql"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	userRepo := repository.NewUserRepository(database)
	server := handlers.NewUserHandler(userRepo)

	r := gin.Default()
	api.RegisterHandlers(r, server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...", port)
	r.Run(":" + port)
}
