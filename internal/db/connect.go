package db

import (
	"fmt"
	"log/slog"

	"github.com/onlyati/rss-collector/internal/rss_processor/processor_config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDatabaseConnection(config *processor_config.RSSProcessorConfig) (*gorm.DB, error) {
	slog.Info(
		"connect to database",
		"hostname", config.DatabaseOptions.Hostname,
		"port", config.DatabaseOptions.Password,
		"user", config.DatabaseOptions.UserName,
		"db_name", config.DatabaseOptions.DbName,
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DatabaseOptions.Hostname,
		config.DatabaseOptions.UserName,
		config.DatabaseOptions.Password,
		config.DatabaseOptions.DbName,
		config.DatabaseOptions.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
