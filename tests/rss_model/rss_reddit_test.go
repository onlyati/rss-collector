package rss_model

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/onlyati/rss-collector/internal/rss_model"
)

func TestRedditRSSUnmarshal(t *testing.T) {
	xmlData := `
	<feed>
		<category term="ansible" label="/r/ansible"/>
		<title>Ansible Subreddit</title>
		<entry>
			<author>
				<name>user123</name>
			</author>
			<category term="ansible" label="/r/ansible"/>
			<content type="html">
				<![CDATA[<div class="md"><!-- SC_OFF -->This is a post content.<!-- SC_ON --></div>]]>
			</content>
			<link href="https://www.reddit.com/r/ansible/comments/example"/>
			<published>2024-12-29T12:00:00Z</published>
			<title>Example Post Title</title>
		</entry>
	</feed>
	`

	var redditRSS rss_model.RedditRSS
	err := xml.Unmarshal([]byte(xmlData), &redditRSS)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Validate top-level fields
	if redditRSS.Category.Term != "ansible" {
		t.Errorf("Expected category term 'ansible', got '%s'", redditRSS.Category.Term)
	}
	if redditRSS.Title != "Ansible Subreddit" {
		t.Errorf("Expected title 'Ansible Subreddit', got '%s'", redditRSS.Title)
	}

	// Validate entry fields
	if len(redditRSS.Entry) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(redditRSS.Entry))
	}

	entry := redditRSS.Entry[0]
	if entry.Author != "user123" {
		t.Errorf("Expected author 'user123', got '%s'", entry.Author)
	}
	if entry.Title != "Example Post Title" {
		t.Errorf("Expected title 'Example Post Title', got '%s'", entry.Title)
	}
	if entry.Link != "https://www.reddit.com/r/ansible/comments/example" {
		t.Errorf("Expected link 'https://www.reddit.com/r/ansible/comments/example', got '%s'", entry.Link)
	}
	expectedContent := "This is a post content."
	if string(entry.Content) != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, entry.Content)
	}
}

func TestRedditEntryContentUnmarshalXML(t *testing.T) {
	xmlData := `
	<content type="html">
		<![CDATA[<div class="md"><!-- SC_OFF -->This is a test content.<!-- SC_ON --></div>]]>
	</content>
	`

	var content rss_model.RedditEntryContent
	err := xml.Unmarshal([]byte(xmlData), &content)
	if err != nil {
		t.Fatalf("Failed to unmarshal content: %v", err)
	}

	expected := "This is a test content."
	if string(content) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(content))
	}
}

func TestRedditRSSCreateRSS(t *testing.T) {
	xmlData := `
	<feed>
		<category term="ansible" label="/r/ansible"/>
		<title>Ansible Subreddit</title>
		<entry>
			<author>
				<name>user123</name>
			</author>
			<category term="ansible" label="/r/ansible"/>
			<content type="html">
				<![CDATA[<div class="md"><!-- SC_OFF -->This is a post content.<!-- SC_ON --></div>]]>
			</content>
			<link href="https://www.reddit.com/r/ansible/comments/example"/>
			<published>2024-12-29T12:00:00Z</published>
			<title>Example Post Title</title>
		</entry>
	</feed>
	`

	var redditRSS rss_model.RedditRSS
	err := xml.Unmarshal([]byte(xmlData), &redditRSS)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	rss, err := redditRSS.CreateRSS()
	if err != nil {
		t.Fatalf("Failed to create RSS: %v", err)
	}

	if rss.Title != "Ansible Subreddit" {
		t.Errorf("Expected RSS title 'Ansible Subreddit', got '%s'", rss.Title)
	}
	if len(rss.Items) != 1 {
		t.Fatalf("Expected 1 RSS item, got %d", len(rss.Items))
	}

	item := rss.Items[0]
	if item.Title != "Example Post Title" {
		t.Errorf("Expected item title 'Example Post Title', got '%s'", item.Title)
	}
	if item.PubDate.Format(time.RFC3339) != "2024-12-29T12:00:00Z" {
		t.Errorf("Expected PubDate '2024-12-29T12:00:00Z', got '%s'", item.PubDate.Format(time.RFC3339))
	}
	if item.Description != "This is a post content." {
		t.Errorf("Expected description 'This is a post content.', got '%s'", item.Description)
	}
	if item.Author != "user123" {
		t.Errorf("Expected author 'user123', got '%s'", item.Author)
	}
	if item.Link != "https://www.reddit.com/r/ansible/comments/example" {
		t.Errorf("Expected link 'https://www.reddit.com/r/ansible/comments/example', got '%s'", item.Link)
	}
	if len(item.Category) != 3 || item.Category[1] != "/r/ansible" {
		t.Errorf("Expected categories ['reddit', '/r/ansible', 'ansible'], got %v", item.Category)
	}
}

func TestRedditRSSGetKafkaKey(t *testing.T) {
	redditRSS := rss_model.RedditRSS{
		Category: rss_model.RedditFeedCategory{
			Term: "ansible",
		},
	}

	kafkaKey := redditRSS.GetKafkaKey()
	expectedKey := "reddit_ansible"
	if kafkaKey != expectedKey {
		t.Errorf("Expected Kafka key '%s', got '%s'", expectedKey, kafkaKey)
	}
}
