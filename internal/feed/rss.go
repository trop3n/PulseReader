package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/jasonkneen/pulsereader/internal/storage"
)

// FetchResult holds the result of fetching a feed.
type FetchResult struct {
	Feed        *storage.Feed
	NewArticles int
	Err         error
}

// FetchFeed fetches a feed URL and stores new articles in the database.
func FetchFeed(db *storage.DB, feedID int64) FetchResult {
	storedFeed, err := db.GetFeed(feedID)
	if err != nil {
		return FetchResult{Err: fmt.Errorf("get feed: %w", err)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fp := gofeed.NewParser()
	parsed, err := fp.ParseURLWithContext(storedFeed.URL, ctx)
	if err != nil {
		return FetchResult{Feed: storedFeed, Err: fmt.Errorf("fetch %s: %w", storedFeed.URL, err)}
	}

	// Update feed metadata if changed
	if parsed.Title != "" && parsed.Title != storedFeed.Title {
		storedFeed.Title = parsed.Title
	}
	if parsed.Link != "" {
		storedFeed.SiteURL = parsed.Link
	}
	if parsed.Description != "" {
		storedFeed.Description = parsed.Description
	}
	if parsed.Image != nil && parsed.Image.URL != "" {
		storedFeed.ImageURL = parsed.Image.URL
	}
	_ = db.AddFeed(storedFeed) // upsert metadata

	newCount := 0
	for _, item := range parsed.Items {
		article := itemToArticle(storedFeed.ID, item)
		err := db.InsertArticle(article)
		if err == nil {
			newCount++
		}
	}

	_ = db.UpdateFeedFetched(storedFeed.ID)

	return FetchResult{Feed: storedFeed, NewArticles: newCount}
}

// AddAndFetchFeed adds a new feed URL, fetches it, and returns the result.
func AddAndFetchFeed(db *storage.DB, url string) FetchResult {
	// First, try to parse the feed to get metadata
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fp := gofeed.NewParser()
	parsed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return FetchResult{Err: fmt.Errorf("fetch %s: %w", url, err)}
	}

	title := parsed.Title
	if title == "" {
		title = url
	}

	storedFeed := &storage.Feed{
		Title:       title,
		URL:         url,
		SiteURL:     parsed.Link,
		Description: parsed.Description,
	}
	if parsed.Image != nil {
		storedFeed.ImageURL = parsed.Image.URL
	}

	if err := db.AddFeed(storedFeed); err != nil {
		return FetchResult{Err: fmt.Errorf("save feed: %w", err)}
	}

	// Re-fetch to get the ID (in case of upsert)
	feeds, _ := db.ListFeeds()
	for _, f := range feeds {
		if f.URL == url {
			storedFeed.ID = f.ID
			break
		}
	}

	newCount := 0
	for _, item := range parsed.Items {
		article := itemToArticle(storedFeed.ID, item)
		if err := db.InsertArticle(article); err == nil {
			newCount++
		}
	}

	_ = db.UpdateFeedFetched(storedFeed.ID)

	return FetchResult{Feed: storedFeed, NewArticles: newCount}
}

// FetchAllFeeds fetches all subscribed feeds.
func FetchAllFeeds(db *storage.DB) []FetchResult {
	feeds, err := db.ListFeeds()
	if err != nil {
		return []FetchResult{{Err: err}}
	}

	results := make([]FetchResult, 0, len(feeds))
	for _, f := range feeds {
		result := FetchFeed(db, f.ID)
		results = append(results, result)
	}
	return results
}

func itemToArticle(feedID int64, item *gofeed.Item) *storage.Article {
	guid := item.GUID
	if guid == "" {
		guid = item.Link
	}
	if guid == "" {
		guid = item.Title
	}

	author := ""
	if item.Author != nil {
		author = item.Author.Name
	}

	// Prefer full content, fall back to description
	content := item.Content
	summary := item.Description

	var publishedAt *time.Time
	if item.PublishedParsed != nil {
		publishedAt = item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		publishedAt = item.UpdatedParsed
	}

	return &storage.Article{
		FeedID:      feedID,
		GUID:        guid,
		Title:       item.Title,
		URL:         item.Link,
		Author:      author,
		Content:     content,
		Summary:     summary,
		PublishedAt: publishedAt,
	}
}
