package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pkg/middleware"
	"pkg/router"
	"syscall"
	"time"

	"business-service/internal/config"
	"business-service/internal/db"
	"business-service/internal/handler"
	"business-service/internal/service"
	"pkg/common_handler"
	"pkg/logger"
)

const (
	serviceName   = "business-service"
	version       = "v1"
	handlerPrefix = "business"
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
	productService := service.NewProductService(dbConn.Queries)

	// Initialize handlers
	healthHandler := common_handler.NewHealthHandler(serviceName, version)
	productHandler := handler.NewProductHandler(productService)

	// Setup routes
	pathBuilder := router.NewPathBuilder(version, handlerPrefix)
	mux := http.NewServeMux()
	mux.HandleFunc(healthHandler.Path(), healthHandler.Health)
	mux.HandleFunc(pathBuilder.Path("products"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Check if it's a price range query
			if r.URL.Query().Get("min_price") != "" && r.URL.Query().Get("max_price") != "" {
				productHandler.GetProductsByPriceRange(w, r)
			} else {
				productHandler.GetProducts(w, r)
			}
		case "POST":
			productHandler.CreateProduct(w, r)
		case "PUT":
			productHandler.UpdateProduct(w, r)
		case "DELETE":
			productHandler.DeleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Setup server
	muxHandler := middleware.LoggingMiddleware(log)(mux)
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
