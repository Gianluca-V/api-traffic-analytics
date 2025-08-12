package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-traffic-analytics/cmd/traffic-ingestor/internal/handler"
	"api-traffic-analytics/cmd/traffic-ingestor/internal/repository"
	"api-traffic-analytics/cmd/traffic-ingestor/internal/service"
	"api-traffic-analytics/internal/pkg/kafka"
	"api-traffic-analytics/internal/pkg/postgres"
	"api-traffic-analytics/internal/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize PostgreSQL
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// Initialize Redis
	rdb, err := redis.GetRedisClient()
	if err != nil {
		log.Fatalf("Unable to connect to redis: %v\n", err)
	}

	// Initialize Kafka Producer
	producer := kafka.CreateProducer(
		[]string{os.Getenv("KAFKA_BROKERS")},
		"traffic-data",
	)
	defer producer.Close()

	// Create repository, service, and handler
	postgresRepo := postgres.NewTrafficDataRepository(db)
	redisRepo := redis.NewCacheRepository(rdb)
	repo := repository.NewRepository(postgresRepo, redisRepo)
	svc := service.NewService(repo, producer)
	h := handler.NewHandler(svc)

	// Create Gin router
	router := gin.Default()
	router.POST("/traffic", h.ReceiveTrafficData)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
