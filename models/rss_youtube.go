package models

import (
	"encoding/xml"
	"strings"
	"time"
)

// Specification for Youtube RSS feed
type YoutubeRSS struct {
	XMLName xml.Name       `xml:"feed"`
	Title   string         `xml:"title"` // Channel's title
	Entry   []YoutubeEntry `xml:"entry"`
}

type YoutubeEntry struct {
	Author  string     `xml:"author>name"` // Who made the video
	Title   string     `xml:"title"`       // Title of the video
	Link    RSSLink    `xml:"link"`        // Link for the video
	PubDate time.Time  `xml:"published"`   // When was the video published
	Media   MediaGroup `xml:"group"`       // Thumbnail and description for the video
}

type MediaGroup struct {
	Description string         `xml:"description"`
	Thumbnail   MediaThumbnail `xml:"thumbnail"`
}

func (rss *YoutubeRSS) GetKafkaKey() string {
	key := strings.ToLower(rss.Title)
	key = strings.ReplaceAll(key, " ", "_")
	return "youtube_" + key
}

func (rss *YoutubeRSS) CreateRSS() (RSS, error) {
	finalRSS := RSS{}
	finalRSS.Title = rss.Title
	finalRSS.Items = []RSSItem{}

	for _, entry := range rss.Entry {
		newItem := RSSItem{
			Title:       entry.Title,
			PubDate:     entry.PubDate,
			Description: entry.Media.Description,
			Category:    []string{"youtube", finalRSS.Title},
			ImageLink:   &entry.Media.Thumbnail.URL,
			Link:        string(entry.Link),
			Author:      entry.Author,
		}
		finalRSS.Items = append(finalRSS.Items, newItem)
	}

	return finalRSS, nil
}
