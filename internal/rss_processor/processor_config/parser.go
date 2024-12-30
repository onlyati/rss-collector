package processor_config

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"
)

type RSSProcessorConfig struct {
	KafkaOptions    ConsumerOptions `yaml:"kafka"` // Kafka related settings
	DatabaseOptions DatabaseOption  `yaml:"db"`
}

type DatabaseOption struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
}

type ConsumerOptions struct {
	Topic   string `yaml:"topic"`    // Topic name where the message should be read
	Server  string `yaml:"server"`   // Kafka server address
	GroupID string `yaml:"group_id"` // Group ID for Kafka consumer group
}

func NewRSSProcessorConfigFromYAML(content []byte) (*RSSProcessorConfig, error) {
	var processor *RSSProcessorConfig
	err := yaml.Unmarshal(content, &processor)
	if err != nil {
		return nil, err
	}

	errFlag := false
	if processor.KafkaOptions.Server == "" {
		slog.Error("failed to parse config", "reason", "missing kafka>server")
		errFlag = true
	}

	if processor.KafkaOptions.Topic == "" {
		slog.Error("failed to parse config", "reason", "missing kafka>topic")
		errFlag = true
	}

	if processor.KafkaOptions.GroupID == "" {
		slog.Error("failed to parse config", "reason", "missing kafka>group_id")
		errFlag = true
	}

	if processor.DatabaseOptions.Hostname == "" {
		slog.Error("failed to parse config", "reason", "missing db>hostname")
		errFlag = true
	}

	if processor.DatabaseOptions.Port == 0 {
		processor.DatabaseOptions.Port = 5432
	}

	if processor.DatabaseOptions.UserName == "" {
		slog.Error("failed to parse config", "reason", "missing db>user")
		errFlag = true
	}

	if processor.DatabaseOptions.Password == "" {
		slog.Error("failed to parse config", "reason", "missing db>password")
		errFlag = true
	}

	if processor.DatabaseOptions.DbName == "" {
		slog.Error("failed to parse config", "reason", "missing db>db_name")
		errFlag = true
	}

	if errFlag {
		return nil, fmt.Errorf("failed to parse config")
	}

	slog.Info("config successfully been read", "config", processor)

	return processor, nil
}
