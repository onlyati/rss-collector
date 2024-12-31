package api

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type APIConfig struct {
	DatabaseOptions DatabaseOption `yaml:"db"`
	ApiOptions      ApiOptions     `yaml:"api"`
}

type ApiOptions struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

type DatabaseOption struct {
	Hostname     string `yaml:"hostname"`
	Port         int    `yaml:"port"`
	UserName     string `yaml:"user"`
	PasswordPath string `yaml:"password_path"`
	Password     string `yaml:"-"`
	DbName       string `yaml:"db_name"`
}

func newAPIConfigFromYAML(content []byte) (*APIConfig, error) {
	var config *APIConfig
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	slog.Info("config has been read", "config", config)

	pw, err := os.ReadFile(config.DatabaseOptions.PasswordPath)
	if err != nil {
		slog.Error("failed to read password file", "file_path", config.DatabaseOptions.PasswordPath)
		return nil, err
	}
	config.DatabaseOptions.Password = string(pw)

	errFlag := false
	if config.DatabaseOptions.Hostname == "" {
		slog.Error("failed to parse config", "reason", "missing db>hostname")
		errFlag = true
	}

	if config.DatabaseOptions.Port == 0 {
		config.DatabaseOptions.Port = 5432
	}

	if config.DatabaseOptions.UserName == "" {
		slog.Error("failed to parse config", "reason", "missing db>user")
		errFlag = true
	}

	if config.DatabaseOptions.Password == "" {
		slog.Error("failed to parse config", "reason", "missing db>password")
		errFlag = true
	}

	if config.DatabaseOptions.DbName == "" {
		slog.Error("failed to parse config", "reason", "missing db>db_name")
		errFlag = true
	}

	if errFlag {
		return nil, fmt.Errorf("failed to parse config")
	}

	slog.Info("config successfully been read")

	return config, nil
}
