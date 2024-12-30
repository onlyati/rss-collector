package rss_processor

import (
	"encoding/json"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/onlyati/rss-collector/internal/db"
	"github.com/onlyati/rss-collector/internal/kafka_ops"
	"github.com/onlyati/rss-collector/internal/rss_model"
	"github.com/onlyati/rss-collector/internal/rss_processor/processor_config"
	"gorm.io/gorm"
)

type RSSProcessor struct {
	Kafka           kafka_ops.RSSKafkaConsumer
	ProcessorConfig processor_config.RSSProcessorConfig
	DatabaseConn    *gorm.DB
}

func NewRSSProcessor(configContent []byte) (*RSSProcessor, error) {
	slog.Info("create new consumer")
	kafka, opts, err := kafka_ops.NewRSSKafkaConsumer(configContent)
	if err != nil {
		return nil, err
	}

	db, err := db.CreateDatabaseConnection(opts)
	if err != nil {
		return nil, err
	}

	processor := RSSProcessor{
		Kafka:           *kafka,
		ProcessorConfig: *opts,
		DatabaseConn:    db,
	}

	return &processor, nil
}

func (processor *RSSProcessor) Read() error {
	slog.Info("reading on Kafka topic", "topic", processor.ProcessorConfig.KafkaOptions.Topic)
	err := processor.Kafka.Consumer.Subscribe(processor.ProcessorConfig.KafkaOptions.Topic, nil)
	if err != nil {
		return err
	}

	for {
		msg, err := processor.Kafka.Consumer.ReadMessage(-1)
		if err != nil {
			return err
		}
		slog.Info("read message", "offset", msg.TopicPartition.Offset)
		processData(msg, processor.DatabaseConn)
	}
}

func processData(m *kafka.Message, db *gorm.DB) {
	// Convert the received message to generic RSS struct
	var rss rss_model.RSS
	err := json.Unmarshal(m.Value, &rss)
	if err != nil {
		slog.Error("failed to parse message from kafka", "error", err)
		return
	}

	// Validate the struct
	err = rss.Validate()
	if err != nil {
		slog.Error("failed to validate WHOLE RSS struct", "error", err)
		return
	}

	// Add thread for the database if not exists yet
	newRSS := rss_model.RSS{
		Title: rss.Title,
	}

	tx := db.Where("title = ?", newRSS.Title).FirstOrCreate(&newRSS)
	if tx.Error != nil {
		slog.Error("failed to insert to database", "error", tx.Error)
		return
	}

	for _, item := range rss.Items {
		if item.Validate() != nil {
			slog.Error("failed to validate rss item", "item", item)
			continue
		}

		item.RSSID = newRSS.ID
		tx := db.Where("link = ?", item.Link).FirstOrCreate(&item)
		if tx.Error != nil {
			slog.Error("failed to insert to database", "error", tx.Error)
			continue
		}
	}
}
