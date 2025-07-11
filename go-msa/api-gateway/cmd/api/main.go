package main

import (
	"api-gateway/internal/config"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"pkg/common_handler"
	"pkg/logger"
	commonMiddleware "pkg/middleware"
)

const (
	serviceName = "api-gateway"
	version     = "v1"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New("api-gateway", cfg.Log.Level)

	healthHandler := common_handler.NewHealthHandler(serviceName, version)
	gatewayHandler, err := handler.NewGatewayHandler(cfg.Services.AuthServiceURL, cfg.Services.BusinessServiceURL)
	if err != nil {
		log.Error("Failed to create gateway common_handler: %v", err)
		return
	}

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.JWKSURL)

	mux := http.NewServeMux()
	mux.HandleFunc(healthHandler.Path(), healthHandler.Health)
	mux.Handle("/", authMiddleware.Authenticate(gatewayHandler))

	muxHandler := commonMiddleware.LoggingMiddleware(log)(mux)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
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
