package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/onlyati/rss-collector/internal/api"
)

var CLI struct {
	Listen struct {
		Config  string `help:"Path for the config YAML file."`
		JsonLog bool   `help:"Log output format is JSON."`
	} `cmd:"" help:"Stream recorded RSS feeds over a REST API."`
}

func main() {
	ctx := kong.Parse(&CLI)
	if CLI.Listen.JsonLog {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	} else {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	switch ctx.Command() {
	case "listen":
		configYAML, err := os.ReadFile(CLI.Listen.Config)
		if err != nil {
			slog.Error("failed to read config file", "error", err)
			os.Exit(1)
		}

		api, err := api.NewRouter(configYAML)
		if err != nil {
			slog.Error("failed to initialize api", "error", err)
			os.Exit(16)
		}
		api.Listen()
	}
}
