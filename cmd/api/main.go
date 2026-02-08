package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "waya/docs" // IMPORTANT: This triggers the Swagger init
	betaworkos "waya/internal/adapters/external/Betaworkos"
	wayaHandler "waya/internal/adapters/handlers/http"
	"waya/internal/adapters/handlers/http/middlewares"
	"waya/internal/adapters/payments/afriex"
	wayaDB "waya/internal/adapters/storage/db"
	"waya/internal/config"
	"waya/internal/core/services"
)

// @title Waya API (Afriex Orchestrator)
// @version 1.0
// @description The high-performance B2B payment orchestration layer for Afriex.
// @termsOfService http://swagger.io/terms/

// @contact.name Waya Support
// @contact.email dev@waya.finance

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Setup Structured Logging (JSON)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// 2. Load Config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 3. Init Database
	db, err := wayaDB.NewDatabase(cfg.Database)
	if err != nil {
		slog.Error("Database init failed", "error", err)
		os.Exit(1)
	}
	defer db.Conn.Close()

	// 1. Init Adapters
    repo := wayaDB.NewRepository(db)
    afriexClient := afriex.NewClient(cfg.Afriex)

	// --- Init Notifier ---
    notifier := betaworkos.NewNotifier(cfg.Waya)

    // 2. Init Service
    // Note: We pass the standard Logger
    svc := services.NewPayoutService(repo, afriexClient, notifier, slog.Default())

    // 3. Init Handler
    payoutHandler := wayaHandler.NewPayoutHandler(svc)

	// 4. Init Echo
	e := echo.New()
	e.HideBanner = true // Keep logs clean

	// 5. Middleware Stack
	e.Use(middleware.Recover()) // Don't crash on panic
	e.Use(middleware.CORS())    // Allow Frontend access
	e.Use(middleware.RequestID()) 
    
    // Custom Slog Middleware for Echo
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"id", v.RequestID,
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
			)
			return nil
		},
	}))

	// 6. Routes
	api := e.Group("/api/v1")
	// Apply the authentication middleware to the whole API group
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return middlewares.APIKeyAuth(next, cfg.Waya.APIKey)
	})
	
    // Health Check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok", "db": "connected"})
	})

	api.POST("/payouts", payoutHandler.HandleBulkPayout)
	api.GET("/payouts/:batch_id", payoutHandler.GetBatchStatus)
	// WEBHOOK ROUTE (The new feature)
// api.POST("/webhooks/afriex", payoutHandler.HandleAfriexWebhook)

	// Swagger Endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 7. Start Server (Graceful Shutdown)
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("Shutting down server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}