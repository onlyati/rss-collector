package rss_reader

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/onlyati/rss-collector/models"
)

// This function read the data from the source link and decode it as it is specified in T.
func GetProcessedRSS[T models.RSSable](sourceLink string) (T, error) {
	var rss T

	// Fetch the feed
	response, err := http.Get(sourceLink)
	if err != nil {
		return rss, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return rss, fmt.Errorf("http request did not return with 200 but with %d", response.StatusCode)
	}

	// Make the decode thing
	decoder := xml.NewDecoder(response.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return rss, err
	}
	return rss, nil
}

// This function convert T type into a generic RSS and convert it to a JSON string,
// which is sent to Kafka.
func ConvertToJSON[T models.RSSable](rss T) (*[]byte, error) {
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

// This is a generic function which pace the action: fetch from source link, convert it to a unified RSS,
// convert it to JSON string then send the data to Kafka.
func collectData[T models.RSSable](reader *RSSReader, componentType string, items []string, urlStart string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, item := range items {
		slog.Info("read "+componentType, "item", item)
		url := fmt.Sprintf("%s%s", urlStart, item)

		// Convert type to a unified RSS type
		rss, err := GetProcessedRSS[T](url)
		if err != nil {
			slog.Error("failed to fetch RSS", "url", url, "error", err)
			continue
		}

		// Convert RSS struct for JSON string
		rssJSON, err := ConvertToJSON(rss)
		if err != nil {
			slog.Error("failed to convert to json", "url", url, "error", err)
			continue
		}

		// Send stuff to kafka
		err = reader.Kafka.Producer.Produce(&kafka.Message{
			Value:          *rssJSON,
			Key:            []byte(rss.GetKafkaKey()),
			TopicPartition: kafka.TopicPartition{Topic: &reader.Kafka.Topic, Partition: kafka.PartitionAny},
		}, reader.Kafka.DeliverChannel)

		if err != nil {
			slog.Error("failed to send data to Kafka", "url", url, "error", err)
			continue
		}
		slog.Info("read "+componentType, "item", item, "status", "done")
	}
}
