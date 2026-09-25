package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sokosplit/sokosplit-core-api/internal/db"
	"github.com/sokosplit/sokosplit-core-api/internal/events"
	"github.com/sokosplit/sokosplit-core-api/internal/handlers"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err := db.AutoMigrate(conn); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	bus := events.NewBus()
	splitHandler := handlers.NewSplitListHandler(conn, bus)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// OpenAPI spec served as static file — see docs/openapi.yaml
	r.StaticFile("/docs", "./docs/openapi.yaml")

	// Wallet service calls back here to confirm on-chain state changes.
	r.POST("/webhooks/wallet-service", handlers.WalletWebhook(conn))

	api := r.Group("/api/v1")
	api.Use(handlers.AuthMiddleware())
	{
		api.POST("/split-lists", splitHandler.Create)
		api.GET("/split-lists", splitHandler.List)
		api.GET("/split-lists/:id", splitHandler.Get)
		api.POST("/split-lists/:id/release", splitHandler.Release)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("sokosplit-core-api listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
