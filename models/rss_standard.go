package models

import (
	"strings"
	"time"
)

// Specification for general RSS feed
type StandardRSS struct {
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title       string      `xml:"title"`
			Link        string      `xml:"link"`
			Description string      `xml:"description"`
			PubDate     RSSItemDate `xml:"pubDate"`
			Author      string      `xml:"author"`
			Creator     string      `xml:"creator"`
			Category    []string    `xml:"category"`
			Enclosure   *struct {
				Url    string `xml:"url,attr"`
				Length int    `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (rss *StandardRSS) GetKafkaKey() string {
	key := strings.ToLower(rss.Channel.Title)
	key = strings.ReplaceAll(key, " ", "_")
	return "rss_" + key
}

func (rss *StandardRSS) CreateRSS() (RSS, error) {
	finalRSS := RSS{}
	finalRSS.Title = rss.Channel.Title
	finalRSS.Items = []RSSItem{}

	for _, item := range rss.Channel.Items {
		author := item.Author
		if author == "" {
			author = item.Creator
		}
		var imageLink *string = nil
		if item.Enclosure != nil {
			imageLink = &item.Enclosure.Url
		}

		newItem := RSSItem{
			Title:       item.Title,
			PubDate:     time.Time(item.PubDate),
			Description: item.Description,
			Author:      author,
			Category:    append(item.Category, finalRSS.Title),
			Link:        item.Link,
			ImageLink:   imageLink,
		}
		finalRSS.Items = append(finalRSS.Items, newItem)
	}

	return finalRSS, nil
}
