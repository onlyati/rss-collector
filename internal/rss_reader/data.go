package rss_reader

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/onlyati/rss-collector/internal/kafka_ops/events"
	"github.com/onlyati/rss-collector/internal/rss_model"
)

// Collect data from various data source
func (reader *RSSReader) CollectData() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go collectData[*rss_model.YoutubeRSS](
		reader,
		"youtube",
		reader.RSSConfig.YoutubeChannels,
		"https://www.youtube.com/feeds/videos.xml?channel_id=",
		&wg,
	)

	wg.Add(1)
	go collectData[*rss_model.RedditRSS](
		reader,
		"reddit",
		reader.RSSConfig.RedditThreads,
		"https://www.reddit.com/r/",
		&wg,
	)

	wg.Add(1)
	go collectData[*rss_model.StandardRSS](
		reader,
		"standard",
		reader.RSSConfig.StandardLinks,
		"",
		&wg,
	)

	wg.Add(1)
	go collectData[*rss_model.CrunchyrollRSS](
		reader,
		"crunchyroll",
		[]string{"https://feeds.feedburner.com/crunchyroll/rss/anime"},
		"",
		&wg,
	)

	wg.Wait()
}

// This is a generic function which pace the action: fetch from source link, convert it to a unified RSS,
// convert it to JSON string then send the data to Kafka.
func collectData[T rss_model.RSSable](reader *RSSReader, componentType string, items []string, urlStart string, wg *sync.WaitGroup) {
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
