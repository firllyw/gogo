package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/firllyw/gogo/internal/api/handlers"
	"github.com/firllyw/gogo/internal/api/middleware"
)

// SetupRouter configures the API routes
func SetupRouter(
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Add middleware for request logging, recovery, CORS, etc.
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Login</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					background-color: #f5f5f5;
				}
				.login-container {
					background-color: white;
					padding: 2rem;
					border-radius: 5px;
					box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
					text-align: center;
				}
				.btn {
					background-color: #4285f4;
					color: white;
					border: none;
					padding: 10px 20px;
					border-radius: 4px;
					cursor: pointer;
					font-size: 16px;
					text-decoration: none;
					display: inline-block;
					margin-top: 20px;
				}
				.btn:hover {
					background-color: #3367d6;
				}
			</style>
		</head>
		<body>
			<div class="login-container">
				<h1>Welcome to User Service</h1>
				<p>Please login to continue</p>
				<a href="/api/v1/auth/login" class="btn">Login with OAuth</a>
			</div>
		</body>
		</html>
		`
		c.Header("Content-Type", "text/html")
		c.String(200, html)
	})


	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			auth.GET("/login", authMiddleware.InitiateLogin)
			auth.GET("/callback", authMiddleware.HandleCallback)
		}

		// User endpoints
		users := v1.Group("/users")
		{
			// Public endpoints
			users.POST("", userHandler.Create)
			
			// Protected endpoints with OAuth auth
			users.Use(authMiddleware.Authenticate()) // comment this line to disable OAuth
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	return r
}