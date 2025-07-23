package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"context"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"os/signal"
	"pkg/common_handler"
	"pkg/logger"
	"pkg/middleware"
	"pkg/router"
	"syscall"
	"time"
)

const (
	serviceName   = "auth-service"
	version       = "v1"
	handlerPrefix = "auth"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log := logger.New(serviceName, cfg.Log.Level)

	// Connect to database
	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		log.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	// Initialize services
	userRepo := repository.NewPostgresUserRepository(db.DB)
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	healthHandler := common_handler.NewHealthHandler(serviceName, version)
	authHandler := handler.NewAuthHandler(authService)

	// Setup routes
	pathBuilder := router.NewPathBuilder(version, handlerPrefix)
	mux := http.NewServeMux()
	mux.HandleFunc(healthHandler.Path(), healthHandler.Health)
	mux.HandleFunc(pathBuilder.Path("login"), authHandler.Login)
	mux.HandleFunc(pathBuilder.Path("register"), authHandler.Register)

	// Setup server
	muxHandler := middleware.LoggingMiddleware(log)(mux)
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: muxHandler,
	}

	// Start server in a goroutine
	go func() {
		log.Info("%s (SQLC + PGX) starting on %s", serviceName, addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Info("Server exited")
}
