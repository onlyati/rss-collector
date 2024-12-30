package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/onlyati/rss-collector/internal/db"
	"github.com/onlyati/rss-collector/internal/rss_processor"
	"github.com/onlyati/rss-collector/internal/rss_processor/processor_config"
)

var CLI struct {
	Process struct {
		Config  string `help:"Path for the config YAML file."`
		JsonLog bool   `help:"Log output format is JSON."`
	} `cmd:"" help:"Process RSS feeds from Kafka."`
	DbMigration struct {
		Config  string `help:"Path for the config YAML file."`
		JsonLog bool   `help:"Log output format is JSON."`
	} `cmd:"" help:"Perform database migration"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if CLI.Process.JsonLog {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	} else {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	switch ctx.Command() {
	case "process":
		config, err := os.ReadFile(CLI.Process.Config)
		if err != nil {
			slog.Error("failed to read config file", "error", err)
			os.Exit(1)
		}

		processor, err := rss_processor.NewRSSProcessor(config)
		if err != nil {
			slog.Error("failed to read RSS feed", "error", err)
		} else {
			processor.Read()

		}
	case "db-migration":
		config, err := os.ReadFile(CLI.DbMigration.Config)
		if err != nil {
			slog.Error("failed to read config file", "error", err)
			os.Exit(1)
		}

		parsedConfig, err := processor_config.NewRSSProcessorConfigFromYAML(config)
		if err != nil {
			slog.Error("failed to parse config", "error", err)
			os.Exit(1)
		}

		db.DatabaseAutoMigration(parsedConfig)
	}
}
