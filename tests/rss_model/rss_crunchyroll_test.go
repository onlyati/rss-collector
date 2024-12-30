package rss_model

import (
	"encoding/xml"
	"testing"

	"github.com/onlyati/rss-collector/internal/rss_model"
)

// Test XML unmarshaling for CrunchyrollRSS
func TestCrunchyrollRSSUnmarshal(t *testing.T) {
	xmlData := `
	<rss>
		<channel>
			<title>Crunchyroll Anime</title>
			<item>
				<title>Naruto (Dub)</title>
				<link>https://example.com/naruto</link>
				<description><![CDATA[<img src="https://example.com/image.jpg" /><br />Episode 1 Description]]></description>
				<pubDate>Wed, 29 Dec 2024 15:04:05 +0000</pubDate>
				<category>Action</category>
				<thumbnail url="https://example.com/image1.jpg" width="1280" height="720" />
				<thumbnail url="https://example.com/image2.jpg" width="1920" height="1080" />
				<publisher>Crunchyroll</publisher>
			</item>
			<item>
				<title>One Piece</title>
				<link>https://example.com/onepiece</link>
				<description><![CDATA[<img src="https://example.com/image.jpg" /><br />Episode 1000 Description]]></description>
				<pubDate>Wed, 27 Dec 2024 15:04:05 +0000</pubDate>
				<category>Adventure</category>
				<thumbnail url="https://example.com/image3.jpg" width="640" height="360" />
				<thumbnail url="https://example.com/image4.jpg" width="1280" height="720" />
				<publisher>Crunchyroll</publisher>
			</item>
		</channel>
	</rss>
	`

	var crunchyrollRSS rss_model.CrunchyrollRSS
	err := xml.Unmarshal([]byte(xmlData), &crunchyrollRSS)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	if crunchyrollRSS.Channel.Title != "Crunchyroll Anime" {
		t.Errorf("Expected title 'Crunchyroll Anime', got '%s'", crunchyrollRSS.Channel.Title)
	}

	if len(crunchyrollRSS.Channel.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(crunchyrollRSS.Channel.Items))
	}

	item := crunchyrollRSS.Channel.Items[0]
	if item.Title != "Naruto (Dub)" {
		t.Errorf("Expected first item title 'Naruto (Dub)', got '%s'", item.Title)
	}
	if item.Link != "https://example.com/naruto" {
		t.Errorf("Expected link 'https://example.com/naruto', got '%s'", item.Link)
	}
	if item.Description != "Episode 1 Description" {
		t.Errorf("Expected description 'Episode 1 Description', got '%s'", item.Description)
	}
	if len(item.Thumbnail) != 2 {
		t.Errorf("Expected 2 thumbnails, got %d", len(item.Thumbnail))
	}
	if item.Thumbnail[1].URL != "https://example.com/image2.jpg" {
		t.Errorf("Expected second thumbnail URL 'https://example.com/image2.jpg', got '%s'", item.Thumbnail[1].URL)
	}
}

func TestCrunchyrollRSSCreateRSS(t *testing.T) {
	xmlData := `
	<rss>
		<channel>
			<title>Crunchyroll Anime</title>
			<item>
				<title>Naruto (Dub)</title>
				<link>https://example.com/naruto</link>
				<description><![CDATA[<img src="https://example.com/image.jpg" /><br />Episode 1 Description]]></description>
				<pubDate>Wed, 29 Dec 2024 15:04:05 +0000</pubDate>
				<category>Action</category>
				<thumbnail url="https://example.com/image1.jpg" width="1280" height="720" />
				<publisher>Crunchyroll</publisher>
			</item>
			<item>
				<title>One Piece</title>
				<link>https://example.com/onepiece</link>
				<description><![CDATA[<img src="https://example.com/image.jpg" /><br />Episode 1000 Description]]></description>
				<pubDate>Wed, 27 Dec 2024 15:04:05 +0000</pubDate>
				<category>Adventure</category>
				<thumbnail url="https://example.com/image3.jpg" width="640" height="360" />
				<thumbnail url="https://example.com/image4.jpg" width="1280" height="720" />
				<publisher>Crunchyroll</publisher>
			</item>
		</channel>
	</rss>
	`

	var crunchyrollRSS rss_model.CrunchyrollRSS
	err := xml.Unmarshal([]byte(xmlData), &crunchyrollRSS)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	genericRSS, err := crunchyrollRSS.CreateRSS()
	if err != nil {
		t.Fatalf("Failed to create RSS: %v", err)
	}

	if genericRSS.Title != "Crunchyroll Anime" {
		t.Errorf("Expected RSS title 'Crunchyroll Anime', got '%s'", genericRSS.Title)
	}

	if len(genericRSS.Items) != 1 { // Expects only 1 item after filtering out "Dub"
		t.Errorf("Expected 1 item after filtering, got %d", len(genericRSS.Items))
	}

	item := genericRSS.Items[0]
	if item.Title != "One Piece" {
		t.Errorf("Expected item title 'One Piece', got '%s'", item.Title)
	}
	if item.ImageLink == nil || *item.ImageLink != "https://example.com/image4.jpg" {
		t.Errorf("Expected image link 'https://example.com/image4.jpg', got '%s'", *item.ImageLink)
	}
	if len(item.Category) != 2 || item.Category[1] != "crunchyroll" {
		t.Errorf("Expected categories ['Adventure', 'crunchyroll'], got %v", item.Category)
	}
}

func TestCrunchyrollRSSGetKafkaKey(t *testing.T) {
	var crunchyrollRSS rss_model.CrunchyrollRSS
	kafkaKey := crunchyrollRSS.GetKafkaKey()

	if kafkaKey != "crunchyroll" {
		t.Errorf("Expected Kafka key 'crunchyroll', got '%s'", kafkaKey)
	}
}
