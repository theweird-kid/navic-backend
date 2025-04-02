package main

import (
	"context"
	"fmt"
	"log"
	"navic-backend/internal/database"
	"navic-backend/internal/handlers"
	message_queue "navic-backend/internal/message-queue"
	"navic-backend/internal/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Init Database
	client, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = message_queue.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %s", err)
	}
	defer message_queue.CloseRabbitMQ()

	// API group
	api := r.Group("/api")
	{
		// User Routes
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.Register)

		// Public route for updating location data in real-time
		api.PUT("/devices/:deviceId/location", handlers.UpdateDeviceLocation)
		api.DELETE("/devices/:deviceId/location", handlers.ClearDeviceLocation)

		// Protected routes
		protected := api.Group("/")
		protected.Use(utils.AuthMiddleware())
		{
			protected.POST("/devices", handlers.AddDevice)
			protected.PUT("/devices/:deviceId", handlers.UpdateDevice)
			protected.DELETE("/devices/:deviceId", handlers.DeleteDevice)
			protected.GET("/devices", handlers.GetDevices)
			protected.GET("/devices/:deviceId", handlers.GetDeviceByID)
			protected.GET("/devices/:deviceId/history", handlers.GetDeviceHistory)

			// Send Message
			protected.POST("/devices/:deviceId/message", handlers.SendMessageToDevice)
		}

		// Server Check
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// DB Check
		api.GET("/health", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			err := client.Ping(ctx, nil)
			if err != nil {
				c.JSON(500, gin.H{
					"status": "unhealthy",
					"error":  err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"status": "healthy",
			})
		})
	}

	r.Run("0.0.0.0:8080")
}
