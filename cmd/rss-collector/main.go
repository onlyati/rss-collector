package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/onlyati/rss-collector/internal/rss_reader"
)

var CLI struct {
	Collect struct {
		Config  string `help:"Path for the config YAML file."`
		JsonLog bool   `help:"Log output format is JSON."`
	} `cmd:"" help:"Collect RSS feeds and push to Kafka based on configuration."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "collect":
		if CLI.Collect.JsonLog {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			slog.SetDefault(logger)
		} else {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			slog.SetDefault(logger)
		}

		config, err := os.ReadFile(CLI.Collect.Config)
		if err != nil {
			slog.Error("failed to read config file", "error", err)
			os.Exit(1)
		}

		reader, err := rss_reader.NewRSSReader(config)
		// reader, err := rss_reader.NewRSSReaderFromYAML(config)
		if err != nil {
			slog.Error("failed to read RSS feed", "error", err)
		} else {
			for {
				reader.CollectData()
				reader.Sleep()
			}
		}
	}
}
