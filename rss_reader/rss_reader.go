package rss_reader

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/onlyati/rss-collector/models"
	"gopkg.in/yaml.v3"
)

type RSSReader struct {
	YoutubeChannels []string `yaml:"youtube"`  // List about youtube channels that must be watched
	RedditThreads   []string `yaml:"reddit"`   // List about subreddit that must be watched
	StandardLinks   []string `yaml:"standard"` // General standard RSS links
	CrunchyrollRSS  bool     `yaml:"crunchyroll"`
	Kafka           struct {
		Topic          string           `yaml:"topic"`          // Topic name where message are sent
		Server         string           `yaml:"server"`         // Server connection where producer is connecting
		ConsumerGroup  string           `yaml:"consumer_group"` // Name of the consumer group that receives the data
		Producer       *kafka.Producer  `yaml:"-"`              // Kafka producer instance to be used by goroutines
		DeliverChannel chan kafka.Event `yaml:"-"`              // Delivery channel for Kafka producer
	} `yaml:"kafka"`
}

// This function create a new RSSReader struct based on YAML input content.
// This function also creates a Kafka producer and go routine that checking the delivery reports from Kafka.
func NewRSSReaderFromYAML(content []byte) (*RSSReader, error) {
	var reader *RSSReader
	err := yaml.Unmarshal(content, &reader)
	if err != nil {
		return nil, err
	}

	// Create Kafka producer instance
	err = CreateNewProducer(reader)
	if err != nil {
		return nil, err
	}

	// Return with the initialized structure
	return reader, nil
}

// Collect data from various data source
func (reader *RSSReader) CollectData() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go collectData[*models.YoutubeRSS](
		reader,
		"youtube",
		reader.YoutubeChannels,
		"https://www.youtube.com/feeds/videos.xml?channel_id=",
		&wg,
	)

	wg.Add(1)
	go collectData[*models.RedditRSS](
		reader,
		"reddit",
		reader.RedditThreads,
		"https://www.reddit.com/r/",
		&wg,
	)

	wg.Add(1)
	go collectData[*models.StandardRSS](
		reader,
		"standard",
		reader.StandardLinks,
		"",
		&wg,
	)

	wg.Add(1)
	go collectData[*models.CrunchyrollRSS](
		reader,
		"crunchyroll",
		[]string{"https://feeds.feedburner.com/crunchyroll/rss/anime"},
		"",
		&wg,
	)

	wg.Wait()
}
