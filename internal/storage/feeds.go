package storage

import (
	"database/sql"
	"time"
)

// Feed represents an RSS/Atom feed subscription.
type Feed struct {
	ID          int64
	Title       string
	URL         string
	SiteURL     string
	Description string
	ImageURL    string
	LastFetched *time.Time
	AddedAt     time.Time

	// Computed
	UnreadCount int
	TotalCount  int
}

// AddFeed adds a new feed subscription.
func (db *DB) AddFeed(feed *Feed) error {
	res, err := db.conn.Exec(
		`INSERT INTO feeds (title, url, site_url, description, image_url)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(url) DO UPDATE SET
		   title=excluded.title, site_url=excluded.site_url,
		   description=excluded.description, image_url=excluded.image_url`,
		feed.Title, feed.URL, feed.SiteURL, feed.Description, feed.ImageURL,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	if id > 0 {
		feed.ID = id
	}
	return nil
}

// ListFeeds returns all feeds with unread counts.
func (db *DB) ListFeeds() ([]Feed, error) {
	rows, err := db.conn.Query(`
		SELECT f.id, f.title, f.url, f.site_url, f.description, f.image_url,
		       f.last_fetched, f.added_at,
		       COUNT(a.id) AS total,
		       COUNT(CASE WHEN a.read = 0 THEN 1 END) AS unread
		FROM feeds f
		LEFT JOIN articles a ON a.feed_id = f.id
		GROUP BY f.id
		ORDER BY f.title COLLATE NOCASE
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var f Feed
		var lastFetched sql.NullTime
		err := rows.Scan(
			&f.ID, &f.Title, &f.URL, &f.SiteURL, &f.Description, &f.ImageURL,
			&lastFetched, &f.AddedAt,
			&f.TotalCount, &f.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		if lastFetched.Valid {
			f.LastFetched = &lastFetched.Time
		}
		feeds = append(feeds, f)
	}
	return feeds, rows.Err()
}

// GetFeed retrieves a feed by ID.
func (db *DB) GetFeed(id int64) (*Feed, error) {
	var f Feed
	var lastFetched sql.NullTime
	err := db.conn.QueryRow(`
		SELECT f.id, f.title, f.url, f.site_url, f.description, f.image_url,
		       f.last_fetched, f.added_at,
		       COUNT(a.id) AS total,
		       COUNT(CASE WHEN a.read = 0 THEN 1 END) AS unread
		FROM feeds f
		LEFT JOIN articles a ON a.feed_id = f.id
		WHERE f.id = ?
		GROUP BY f.id
	`, id).Scan(
		&f.ID, &f.Title, &f.URL, &f.SiteURL, &f.Description, &f.ImageURL,
		&lastFetched, &f.AddedAt,
		&f.TotalCount, &f.UnreadCount,
	)
	if err != nil {
		return nil, err
	}
	if lastFetched.Valid {
		f.LastFetched = &lastFetched.Time
	}
	return &f, nil
}

// RemoveFeed deletes a feed and all its articles.
func (db *DB) RemoveFeed(id int64) error {
	_, err := db.conn.Exec(`DELETE FROM feeds WHERE id = ?`, id)
	return err
}

// UpdateFeedFetched updates the last_fetched timestamp.
func (db *DB) UpdateFeedFetched(id int64) error {
	_, err := db.conn.Exec(`UPDATE feeds SET last_fetched = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

// TotalUnreadCount returns the total unread article count across all feeds.
func (db *DB) TotalUnreadCount() int {
	var count int
	db.conn.QueryRow(`SELECT COUNT(*) FROM articles WHERE read = 0`).Scan(&count)
	return count
}
