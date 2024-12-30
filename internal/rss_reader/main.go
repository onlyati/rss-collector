package rss_reader

import (
	"log/slog"
	"sync"
	"time"

	"github.com/onlyati/rss-collector/internal/kafka_ops"
	"github.com/onlyati/rss-collector/internal/rss_model"
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

func (reader *RSSReader) Start() {
	for {
		reader.collectData()
		reader.sleep()
	}
}

func (reader *RSSReader) sleep() {
	now := time.Now()
	sleepDuration := time.Second * time.Duration(reader.RSSConfig.WaitSeconds)
	nextRun := now.Add(sleepDuration)

	slog.Info("sleep before next cycle", "now", now, "next_run", nextRun)
	time.Sleep(sleepDuration)
}

// Collect data from various data source
func (reader *RSSReader) collectData() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go sendData[*rss_model.YoutubeRSS](
		reader,
		"youtube",
		reader.RSSConfig.YoutubeChannels,
		"https://www.youtube.com/feeds/videos.xml?channel_id=",
		&wg,
	)

	wg.Add(1)
	go sendData[*rss_model.RedditRSS](
		reader,
		"reddit",
		reader.RSSConfig.RedditThreads,
		"https://www.reddit.com/r/",
		&wg,
	)

	wg.Add(1)
	go sendData[*rss_model.StandardRSS](
		reader,
		"standard",
		reader.RSSConfig.StandardLinks,
		"",
		&wg,
	)

	wg.Add(1)
	go sendData[*rss_model.CrunchyrollRSS](
		reader,
		"crunchyroll",
		[]string{"https://feeds.feedburner.com/crunchyroll/rss/anime"},
		"",
		&wg,
	)

	wg.Wait()
}
