package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanapoln/capgo-server/app"
	"github.com/tanapoln/capgo-server/app/db"
	"github.com/tanapoln/capgo-server/cmd/server/otel"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("Connecting to database...")
	if err := db.InitDB(ctx); err != nil {
		slog.Error("Error init db", "error", err)
		os.Exit(1)
	}

	router := app.InitRouter()

	shutdownOtel, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		slog.Error("Error setup otel", "error", err)
		os.Exit(1)
	}
	defer shutdownOtel(context.Background())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		slog.Info("Start listening server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error listening", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("Shutting down server...")
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
