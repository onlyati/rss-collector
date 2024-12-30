package rss_model

import (
	"encoding/xml"
	"html"
	"strings"
	"time"
)

// Specification for Reddit RSS feed
type RedditRSS struct {
	XMLName  xml.Name           `xml:"feed"`
	Category RedditFeedCategory `xml:"category"` // Indicate the subreddit category
	Title    string             `xml:"title"`    // Subreddit title
	Entry    []RedditEntry      `xml:"entry"`    // Posts in the subreddit
}

type RedditFeedCategory struct {
	Term  string `xml:"term,attr"`  // It is "ansible" for ansible subreddit
	Label string `xml:"label,attr"` // It is "/r/ansible" for ansible subreddit
}

type RedditEntry struct {
	Author   string             `xml:"author>name"` // User name who made post
	Category RedditFeedCategory `xml:"category"`    // Same category then the main element
	Content  RedditEntryContent `xml:"content"`     // Content of the post
	Link     RSSLink            `xml:"link"`        // Link for the post
	PubDate  time.Time          `xml:"published"`   // When the post was created
	Title    string             `xml:"title"`       // Title of the post
}

// Reddit deliver the content in HTML format, but all character are escaped.
// So they have to converted back and the "frame div" that is put among the content is also removed.
type RedditEntryContent string

func (c *RedditEntryContent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Type               string `xml:"type,attr"`
		RedditEntryContent string `xml:",chardata"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	if v.Type == "html" {
		v.RedditEntryContent = html.UnescapeString(v.RedditEntryContent)
		if strings.Contains(v.RedditEntryContent, "<!-- SC_OFF -->") {
			startPos := strings.Index(v.RedditEntryContent, "<!-- SC_OFF -->")
			endPos := strings.Index(v.RedditEntryContent, "<!-- SC_ON -->")
			v.RedditEntryContent = v.RedditEntryContent[startPos+len("<!-- SC_OFF -->") : endPos]
		}

		if v.RedditEntryContent[0:len("<div class=\"md\">")] == "<div class=\"md\">" {
			v.RedditEntryContent = v.RedditEntryContent[len("<div class=\"md\">") : len(v.RedditEntryContent)-len("</div>")]
		}
	}
	*c = RedditEntryContent(v.RedditEntryContent)
	return nil
}

func (rss *RedditRSS) GetKafkaKey() string {
	return "reddit_" + rss.Category.Term
}

func (rss *RedditRSS) CreateRSS() (RSS, error) {
	finalRSS := RSS{}
	finalRSS.Title = rss.Title

	for _, entry := range rss.Entry {
		newItem := RSSItem{
			Title:       entry.Title,
			PubDate:     entry.PubDate,
			Description: string(entry.Content),
			Category:    []string{"reddit", entry.Category.Label, entry.Category.Term},
			ImageLink:   nil,
			Link:        string(entry.Link),
			Author:      entry.Author,
		}
		finalRSS.Items = append(finalRSS.Items, newItem)
	}

	return finalRSS, nil
}
