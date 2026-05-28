// Command server starts the devops-course API server with optional PostgreSQL
// connectivity and graceful shutdown on SIGINT/SIGTERM.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"

	"github.com/omargallob/devops-course/internal/api"
	"github.com/omargallob/devops-course/internal/database"

	// Enforces build constraints: only compiles on linux/{amd64,arm64} and darwin/{amd64,arm64}.
	_ "github.com/omargallob/devops-course/internal/platform"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Database connection (optional: if DATABASE_URL is not set, run without DB).
	var db *gorm.DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		var err error
		db, err = database.Open(dsn, logger)
		if err != nil {
			slog.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		slog.Info("database ready")
	} else {
		slog.Warn("DATABASE_URL not set, running without database")
	}

	router := api.NewRouter(logger, db)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced shutdown", "error", err)
	}

	fmt.Println("server stopped")
}
