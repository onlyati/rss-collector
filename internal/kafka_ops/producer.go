package kafka_ops

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/onlyati/rss-collector/internal/rss_reader/reader_config"
)

type RSSKafkaProducer struct {
	Producer       *kafka.Producer  // Kafka producer instance to be used by goroutines
	DeliverChannel chan kafka.Event // Delivery channel for Kafka producer
}

func NewRSSKafkaProducer(configContent []byte) (*RSSKafkaProducer, *reader_config.RSSReaderConfig, error) {
	rssConfig, err := reader_config.NewRSSReaderConfigFromYAML(configContent)
	if err != nil {
		return nil, nil, err
	}

	producer, channel, err := createNewProducer(rssConfig.KafkaOptions.Server)
	if err != nil {
		return nil, nil, err
	}

	kafka := RSSKafkaProducer{
		Producer:       producer,
		DeliverChannel: channel,
	}

	return &kafka, rssConfig, nil
}

// This function create a new Kafka produce and register it into RSSReader struct.
// This function also creates a go routine to check Kafka delivery reports.
//
// Args:
//   - server: address of the Kafka
//
// Returns
//   - Kafka producer (in case of error nil)
//   - Delivery report channel (in case of error nil)
//   - Error message (in case of success nil)
func createNewProducer(server string) (*kafka.Producer, chan kafka.Event, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": server,
	})
	if err != nil {
		return nil, nil, err
	}

	repChan := make(chan kafka.Event)

	go func() {
		for {
			e := <-repChan
			m := e.(*kafka.Message)

			if m.TopicPartition.Error != nil {
				slog.Error(
					"delivery failed to Kafka",
					"error", m.TopicPartition.Error,
					"key", m.Key,
				)
			} else {
				slog.Info(
					"message delivered to Kafka",
					"partition", m.TopicPartition.Partition,
					"offset", m.TopicPartition.Offset,
					"key", m.Key,
				)
			}
		}
	}()

	return producer, repChan, nil
}
