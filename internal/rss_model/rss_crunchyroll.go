package rss_model

import (
	"encoding/xml"
	"strings"
	"time"
)

// Specification for Crunchyroll RSS feed
type CrunchyrollRSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string `xml:"title"` // Website title
		Items []struct {
			Title       string                     `xml:"title"`       // Anime title
			Link        string                     `xml:"link"`        // Link to play anime
			Description CrunchyrollItemDescription `xml:"description"` // Description about the episode
			PubDate     RSSItemDate                `xml:"pubDate"`     // Published date of the episode
			Category    []string                   `xml:"category"`    // Category of the item
			Thumbnail   []MediaThumbnail           `xml:"thumbnail"`   // Crunchyroll provides more thumbnails in different size
			Publisher   string                     `xml:"publisher"`   // Name of the episode's publisher
		} `xml:"item"`
	} `xml:"channel"`
}

// The content is basically a string, but it start with a small thumbnail and linebreak.
// This unmarshal method removed this frame.
type CrunchyrollItemDescription string

func (desc *CrunchyrollItemDescription) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	cutPos := strings.Index(v, "<br />")
	if cutPos > -1 {
		v = v[cutPos+len("<br />"):]
	}
	*desc = CrunchyrollItemDescription(v)
	return nil
}

func (*CrunchyrollRSS) GetKafkaKey() string {
	return "crunchyroll"
}

func (rss *CrunchyrollRSS) CreateRSS() (RSS, error) {
	finalRSS := RSS{}
	finalRSS.Title = rss.Channel.Title
	finalRSS.Items = []RSSItem{}

	for _, item := range rss.Channel.Items {
		// Exclude dub versions
		if strings.Contains(item.Title, "Dub)") {
			continue
		}

		// Find the biggest image
		imageLink := ""
		imageWidth := 0
		for _, thumbnail := range item.Thumbnail {
			if thumbnail.Width > imageWidth {
				imageLink = thumbnail.URL
				imageWidth = thumbnail.Width
			}
		}

		// Create the common RSS item
		newItem := RSSItem{
			Title:       item.Title,
			PubDate:     time.Time(item.PubDate),
			Description: string(item.Description),
			Category:    append(item.Category, "crunchyroll"),
			ImageLink:   &imageLink,
			Link:        item.Link,
			Author:      item.Publisher,
		}
		finalRSS.Items = append(finalRSS.Items, newItem)
	}

	return finalRSS, nil
}
