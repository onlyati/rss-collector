package rss_reader

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// This function create a new Kafka produce and register it into RSSReader struct.
// This function also creates a go routine to check Kafka delivery reports.
func CreateNewProducer(reader *RSSReader) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": reader.Kafka.Server,
	})
	if err != nil {
		return err
	}
	reader.Kafka.Producer = producer

	repChan := make(chan kafka.Event)
	reader.Kafka.DeliverChannel = repChan

	go func() {
		for {
			e := <-reader.Kafka.DeliverChannel
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

	return nil
}
