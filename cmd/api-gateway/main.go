package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"api-traffic-analytics/cmd/api-gateway/internal/config"
	"api-traffic-analytics/cmd/api-gateway/internal/handler"
	"api-traffic-analytics/cmd/api-gateway/internal/middleware"
	"api-traffic-analytics/cmd/api-gateway/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin to release mode in production
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize services
	proxyService := service.NewProxyService(cfg)

	// Initialize handler
	apiHandler := handler.NewHandler(proxyService, cfg)

	// Initialize router
	router := setupRouter(apiHandler, cfg)

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("API Gateway started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API Gateway exited")
}

func setupRouter(handler *handler.Handler, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logging())
	router.Use(middleware.CORS())

	// Public routes (no auth required)
	public := router.Group("/")
	{
		public.GET("/health", handler.HealthCheck)
		public.POST("/traffic", handler.ReceiveTrafficData) // Para compatibilidad
	}

	// Protected routes (auth required)
	protected := router.Group("/")
	protected.Use(middleware.APIKeyAuth(cfg.APIKey))
	protected.Use(middleware.RateLimit())

	{
		// Analytics endpoints
		protected.GET("/analytics", handler.GetAnalytics)
		protected.GET("/analytics/:locationId", handler.GetAnalyticsByLocation)

		// Alerts endpoints
		protected.GET("/alerts", handler.GetAlerts)
		protected.GET("/alerts/:locationId", handler.GetAlertsByLocation)

		// Traffic data endpoints
		protected.GET("/traffic", handler.GetTrafficData)
		protected.GET("/traffic/:locationId", handler.GetTrafficDataByLocation)

		// Proxy to internal services
		protected.Any("/services/*path", handler.ProxyToService)
	}

	return router
}
