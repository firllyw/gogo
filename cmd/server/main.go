package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/firllyw/gogo/config"
	"github.com/firllyw/gogo/internal/api/handlers"
	"github.com/firllyw/gogo/internal/api/middleware"
	"github.com/firllyw/gogo/internal/api/routes"
	"github.com/firllyw/gogo/internal/repository/postgres"
	"github.com/firllyw/gogo/internal/service"
	"github.com/firllyw/gogo/pkg/db"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables or defaults.")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Initialize middleware
	var authMiddleware *middleware.AuthMiddleware
	if cfg.OAuth.Enabled {
		authMiddleware = middleware.NewAuthMiddleware(
			cfg.OAuth.ClientID,
			cfg.OAuth.ClientSecret,
			cfg.OAuth.RedirectURL,
			cfg.OAuth.AuthURL,
			cfg.OAuth.TokenURL,
			cfg.OAuth.UserInfoURL,
			cfg.OAuth.JWTSecret,
		)
	} else {
		log.Println("OAuth is disabled. Using dummy authentication middleware.")
		authMiddleware = middleware.NewAuthMiddleware("", "", "", "", "", "", "dummy-secret")
	}

	// Setup router
	router := routes.SetupRouter(userHandler, authMiddleware)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}