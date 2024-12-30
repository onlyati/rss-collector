package rss_reader

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/onlyati/rss-collector/internal/kafka_ops/events"
	"github.com/onlyati/rss-collector/internal/rss_model"
)

// This is a generic function which pace the action: fetch from source link, convert it to a unified RSS,
// convert it to JSON string then send the data to Kafka.
func sendData[T rss_model.RSSable](reader *RSSReader, componentType string, items []string, urlStart string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, item := range items {
		slog.Info("read "+componentType, "item", item)
		url := fmt.Sprintf("%s%s", urlStart, item)

		// Fetch the feed by a simple HTTP GET request
		response, err := http.Get(url)
		if err != nil {
			slog.Error("failed to fetch RSS", "url", url, "error", err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			slog.Error("http request did not return with 200", "status_code", response.StatusCode)
			continue
		}

		err = events.SendFeedToKafka[T](
			response.Body,
			reader.RSSConfig.KafkaOptions.Topic,
			reader.Kafka.Producer,
			reader.Kafka.DeliverChannel,
		)
		if err != nil {
			slog.Error("failed to send event to kafka", "url", url, "error", err)
			continue
		}
		slog.Info("read "+componentType, "item", item, "status", "done")
	}
}
