package reader_config

import (
	"gopkg.in/yaml.v3"
)

type RSSReaderConfig struct {
	YoutubeChannels []string     `yaml:"youtube"`      // List about youtube channels that must be watched
	RedditThreads   []string     `yaml:"reddit"`       // List about subreddit that must be watched
	StandardLinks   []string     `yaml:"standard"`     // General standard RSS links
	CrunchyrollRSS  bool         `yaml:"crunchyroll"`  // Collect newly released anime
	KafkaOptions    KafkaOptions `yaml:"kafka"`        // Kafka related settings
	WaitSeconds     int          `yaml:"wait_seconds"` // Wait time between two calls
}

type KafkaOptions struct {
	Topic  string `yaml:"topic"`  // Topic name where message are sent
	Server string `yaml:"server"` // Server connection where producer is connecting
}

// This function create a new RSSReader struct based on YAML input content.
// This function also creates a Kafka producer and go routine that checking the delivery reports from Kafka.
func NewRSSReaderConfigFromYAML(content []byte) (*RSSReaderConfig, error) {
	var reader *RSSReaderConfig
	err := yaml.Unmarshal(content, &reader)
	if err != nil {
		return nil, err
	}

	// Return with the initialized structure
	return reader, nil
}
