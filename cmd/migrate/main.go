package main

import (
	"log/slog"
	"os"

	"github.com/tanapoln/capgo-server/app/db"
)

func main() {
	slog.Info("Connecting to database...")
	if err := db.InitDB(); err != nil {
		slog.Error("Error init db", "error", err)
		os.Exit(1)
	}

	slog.Info("Running migration...")
	if err := db.RunMigration(); err != nil {
		slog.Error("Error running migration", "error", err)
		os.Exit(1)
	}

	slog.Info("Migration done")
}
