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

	"github.com/tanapoln/capgo-server/app"
	"github.com/tanapoln/capgo-server/app/db"
	"github.com/tanapoln/capgo-server/cmd/server/otel"
	"github.com/tanapoln/capgo-server/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("Connecting to database...")
	if err := db.InitDB(ctx); err != nil {
		slog.Error("Error init db", "error", err)
		os.Exit(1)
	}

	shutdownOtel, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		slog.Error("Error setup otel", "error", err)
		os.Exit(1)
	}
	defer shutdownOtel(context.Background())

	userSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Get().CapgoUserPort),
		Handler: app.InitRouter(),
	}

	mgmtSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Get().CapgoManagementPort),
		Handler: app.InitMgmtRouter(),
	}

	go func() {
		slog.Info("Start listening server", "address", userSrv.Addr)
		if err := userSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error listening", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		slog.Info("Start listening management server", "address", mgmtSrv.Addr)
		if err := mgmtSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error listening", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("Shutting down server...")
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := mgmtSrv.Shutdown(ctx); err != nil {
		slog.Error("Management server forced to shutdown", "error", err)
	}
	if err := userSrv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
