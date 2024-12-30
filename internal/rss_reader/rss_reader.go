package rss_reader

import (
	"log/slog"
	"time"

	"github.com/onlyati/rss-collector/internal/kafka_ops"
	"github.com/onlyati/rss-collector/internal/rss_reader/reader_config"
)

type RSSReader struct {
	Kafka     kafka_ops.RSSKafkaProducer
	RSSConfig reader_config.RSSReaderConfig
}

func NewRSSReader(configContent []byte) (*RSSReader, error) {
	kafka, config, err := kafka_ops.NewRSSKafkaProducer(configContent)
	if err != nil {
		return nil, err
	}

	reader := RSSReader{
		Kafka:     *kafka,
		RSSConfig: *config,
	}

	return &reader, nil
}

func (reader *RSSReader) Sleep() {
	now := time.Now()
	sleepDuration := time.Second * time.Duration(reader.RSSConfig.WaitSeconds)
	nextRun := now.Add(sleepDuration)

	slog.Info("sleep before next cycle", "now", now, "next_run", nextRun)
	time.Sleep(sleepDuration)
}
