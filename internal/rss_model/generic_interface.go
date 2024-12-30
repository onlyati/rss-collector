package rss_model

import (
	"time"
)

// Generic RSS struct, all other type must be converted into this generic one.
type RSS struct {
	Title string    `json:"title"` // Title of the site
	Items []RSSItem `json:"items"` // Content of the RSS
}

// Generic RSS items.
type RSSItem struct {
	Title       string    `json:"title"`       // Article/Video/Whatever title.
	PubDate     time.Time `json:"published"`   // Published data.
	Description string    `json:"description"` // Description of the item.
	Category    []string  `json:"category"`    // Categories where the item is assigned.
	ImageLink   *string   `json:"image_link"`  // Optional thumbnail image link.
	Link        string    `json:"link"`        // Link for the source.
	Author      string    `json:"author"`      // Who created the source.
}

// This interface must be implemented by all type of RSS structures.
type RSSable interface {
	CreateRSS() (RSS, error) // This method converts the special type RSS to generic RSS type.
	GetKafkaKey() string     // This can help for Kafka to distribute among partitions
}

// Used to store thumbnail images
type MediaThumbnail struct {
	URL    string `xml:"url,attr"`    // Link for the image
	Width  int    `xml:"width,attr"`  // Original width of the image
	Height int    `xml:"height,attr"` // Original height of the image
}
