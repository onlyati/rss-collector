package rss_model

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Generic RSS struct, all other type must be converted into this generic one.
type RSS struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title string    `json:"title" gorm:"type:text;not null;unique"`                    // Title of the site
	Items []RSSItem `json:"items" gorm:"foreignKey:RSSID;constraint:OnDelete:CASCADE"` // Content of the RSS
}

func (r *RSS) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil { // Assign UUIDv7 if not already set
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		r.ID = id
	}
	return nil
}

func (rss *RSS) Validate() error {
	if rss.Title == "" {
		return fmt.Errorf("rss: title missing")
	}

	return nil
}

// Generic RSS items.
type RSSItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;;primaryKey" json:"id"`  // Use UUIDv7
	RSSID       uuid.UUID      `gorm:"type:uuid;not null" json:"rss_id"` // Key for the header item
	Title       string         `gorm:"not null" json:"title"`            // Article/Video/Whatever title.
	PubDate     time.Time      `gorm:"not null;index" json:"published"`  // Published data.
	Description string         `gorm:"" json:"description"`              // Description of the item.
	Category    pq.StringArray `gorm:"type:text[]" json:"category"`      // Categories where the item is assigned.
	ImageLink   *string        `gorm:"type:text" json:"image_link"`      // Optional thumbnail image link.
	Link        string         `gorm:"not null;unique" json:"link"`      // Link for the source.
	Author      string         `gorm:"" json:"author"`                   // Who created the source.
}

func (ri *RSSItem) BeforeCreate(tx *gorm.DB) error {
	if ri.ID == uuid.Nil { // Assign UUIDv7 if not already set
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		ri.ID = id
	}
	return nil
}

func (item *RSSItem) Validate() error {
	errorMsg := []string{}

	if item.Title == "" {
		errorMsg = append(errorMsg, "missing title")
	}

	if item.Link == "" {
		errorMsg = append(errorMsg, "missing link")
	}

	if len(errorMsg) > 0 {
		slog.Error("validation failed on item", "item", item, "error", errorMsg)
		return fmt.Errorf("validation failed on item")
	}

	return nil
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
