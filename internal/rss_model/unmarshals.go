package rss_model

import (
	"encoding/xml"
	"strings"
	"time"
)

// This time is a wrapper to convert from standard RSS date format during XML decoding
type RSSItemDate time.Time

func (date *RSSItemDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	words := strings.Split(v, " ")
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	if words[len(words)-1] == "GMT" {
		layout = "Mon, 02 Jan 2006 15:04:05 GMT"
	}
	formattedDate, err := time.Parse(layout, v)
	if err != nil {
		return err
	}
	*date = RSSItemDate(formattedDate)
	return nil
}

// This is a wrapper that get the link from the internal element's attribute
type RSSLink string

func (l *RSSLink) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Href string `xml:"href,attr"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	*l = RSSLink(v.Href)
	return nil
}
