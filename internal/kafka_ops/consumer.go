package kafka_ops

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/onlyati/rss-collector/internal/rss_processor/processor_config"
)

type RSSKafkaConsumer struct {
	Consumer *kafka.Consumer // Kafka consumer instance to be used
}

func NewRSSKafkaConsumer(contentConfig []byte) (*RSSKafkaConsumer, *processor_config.RSSProcessorConfig, error) {
	config, err := processor_config.NewRSSProcessorConfigFromYAML(contentConfig)
	if err != nil {
		return nil, nil, err
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaOptions.Server,
		"group.id":          config.KafkaOptions.GroupID, // Consumer group ID
		"auto.offset.reset": "earliest",                  // Start reading from the earliest message
	})
	if err != nil {
		return nil, nil, err
	}

	kafka := RSSKafkaConsumer{
		Consumer: consumer,
	}

	return &kafka, config, nil
}
