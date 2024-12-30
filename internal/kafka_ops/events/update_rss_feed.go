package events

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/onlyati/rss-collector/internal/rss_model"
)

// This function read the data from the source link and decode it as it is specified in T.
func convertToGenericRSS[T rss_model.RSSable](xmlContent io.Reader) (T, error) {
	var rss T
	decoder := xml.NewDecoder(xmlContent)
	err := decoder.Decode(&rss)
	if err != nil {
		return rss, err
	}
	return rss, nil
}

// This function convert T type into a generic RSS and convert it to a JSON string,
// which is sent to Kafka.
func convertToJSON[T rss_model.RSSable](rss T) (*[]byte, error) {
	finalRSS, err := rss.CreateRSS()
	if err != nil {
		return nil, err
	}

	rssMessage, err := json.Marshal(finalRSS)
	if err != nil {
		return nil, err
	}
	return &rssMessage, nil
}

func SendFeedToKafka[T rss_model.RSSable](
	xmlContent io.Reader,
	topic string,
	producer *kafka.Producer,
	deliverChan chan kafka.Event,
) error {
	// Convert type to a unified RSS type
	rss, err := convertToGenericRSS[T](xmlContent)
	if err != nil {
		slog.Error("failed to convert to XML", "error", err)
		return err
	}

	// Convert RSS struct for JSON string
	rssJSON, err := convertToJSON(rss)
	if err != nil {
		slog.Error("failed to convert to json", "error", err)
		return err
	}

	// Send stuff to kafka
	err = producer.Produce(&kafka.Message{
		Value:          *rssJSON,
		Key:            []byte(rss.GetKafkaKey()),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	}, deliverChan)
	if err != nil {
		return err
	}

	return nil
}
