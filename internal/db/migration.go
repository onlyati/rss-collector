package db

import (
	"log/slog"

	"github.com/onlyati/rss-collector/internal/rss_model"
	"github.com/onlyati/rss-collector/internal/rss_processor/processor_config"
)

func DatabaseAutoMigration(config *processor_config.RSSProcessorConfig) {
	db, err := CreateDatabaseConnection(config)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}

	// Run migrations
	err = db.AutoMigrate(&rss_model.RSS{}, &rss_model.RSSItem{})
	if err != nil {
		slog.Error("failed to migrate database", "error", err)
	}
	slog.Info("database migration is done")
}
