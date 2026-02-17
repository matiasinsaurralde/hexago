/*
Copyright © 2026
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/padiazg/api/pkg/logger"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the hexago-api-chi HTTP server",
	Long: `Start the hexago-api-chi HTTP API server with graceful shutdown support.

The server will listen for SIGINT (Ctrl+C) and SIGTERM signals
and perform a graceful shutdown with a configurable timeout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()

		// Initialize logger from config
		log := logger.New(&logger.Config{
			Level:  cfg.LogLevel,
			Format: cfg.LogFormat,
		})

		log.Info("Starting hexago-api-chi HTTP server...")

		// Setup router and routes
		router := chi.NewRouter()

		// Middleware
		router.Use(chimiddleware.Logger)
		router.Use(chimiddleware.Recoverer)

		// Setup routes
		setupRoutes(router)

		// Configure server
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		}

		// Channel to capture OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		// Channel for server errors
		errChan := make(chan error, 1)

		// Start server in goroutine
		go func() {
			log.Info("Server listening on port %d", cfg.Server.Port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}()

		// Wait for signal or error
		select {
		case sig := <-sigChan:
			log.Info("Received signal: %v, initiating graceful shutdown...", sig)

			// Create shutdown context with timeout
			shutdownCtx, shutdownCancel := context.WithTimeout(
				context.Background(),
				cfg.Server.ShutdownTimeout,
			)
			defer shutdownCancel()

			// Attempt graceful shutdown
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Error("Server shutdown error: %v", err)
				return err
			}

			log.Info("Server stopped gracefully")
			return nil

		case err := <-errChan:
			log.Error("Server error: %v", err)
			return err
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// setupRoutes configures the HTTP routes
func setupRoutes(router chi.Router) {
	// API health check (lightweight)
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"hexago-api-chi"}`)
	})

	// TODO: Add your routes here
	// Example:
	// router.Route("/api/v1", func(r chi.Router) {
	//     r.Get("/users", handlers.ListUsers)
	//     r.Post("/users", handlers.CreateUser)
	// })
}
