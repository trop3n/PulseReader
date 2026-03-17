package storage

import (
	"database/sql"
	"time"
)

// Article represents a feed article/entry.
type Article struct {
	ID          int64
	FeedID      int64
	GUID        string
	Title       string
	URL         string
	Author      string
	Content     string
	Summary     string
	PublishedAt *time.Time
	FetchedAt   time.Time
	Read        bool
	ReadAt      *time.Time

	// Joined
	FeedTitle string
}

// InsertArticle inserts an article, ignoring duplicates.
func (db *DB) InsertArticle(a *Article) error {
	var publishedAt interface{}
	if a.PublishedAt != nil {
		publishedAt = *a.PublishedAt
	}
	_, err := db.conn.Exec(
		`INSERT OR IGNORE INTO articles (feed_id, guid, title, url, author, content, summary, published_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.FeedID, a.GUID, a.Title, a.URL, a.Author, a.Content, a.Summary, publishedAt,
	)
	return err
}

// ListArticles returns articles for a feed, ordered by publish date descending.
// If feedID is 0, returns articles from all feeds.
// If unreadOnly is true, only returns unread articles.
func (db *DB) ListArticles(feedID int64, unreadOnly bool, limit, offset int) ([]Article, error) {
	query := `
		SELECT a.id, a.feed_id, a.guid, a.title, a.url, a.author,
		       a.summary, a.published_at, a.fetched_at, a.read, a.read_at,
		       f.title AS feed_title
		FROM articles a
		JOIN feeds f ON f.id = a.feed_id
		WHERE 1=1
	`
	args := []interface{}{}

	if feedID > 0 {
		query += ` AND a.feed_id = ?`
		args = append(args, feedID)
	}
	if unreadOnly {
		query += ` AND a.read = 0`
	}

	query += ` ORDER BY COALESCE(a.published_at, a.fetched_at) DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		var publishedAt, readAt sql.NullTime
		err := rows.Scan(
			&a.ID, &a.FeedID, &a.GUID, &a.Title, &a.URL, &a.Author,
			&a.Summary, &publishedAt, &a.FetchedAt, &a.Read, &readAt,
			&a.FeedTitle,
		)
		if err != nil {
			return nil, err
		}
		if publishedAt.Valid {
			a.PublishedAt = &publishedAt.Time
		}
		if readAt.Valid {
			a.ReadAt = &readAt.Time
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}

// GetArticle retrieves a full article by ID (including content).
func (db *DB) GetArticle(id int64) (*Article, error) {
	var a Article
	var publishedAt, readAt sql.NullTime
	err := db.conn.QueryRow(`
		SELECT a.id, a.feed_id, a.guid, a.title, a.url, a.author,
		       a.content, a.summary, a.published_at, a.fetched_at, a.read, a.read_at,
		       f.title AS feed_title
		FROM articles a
		JOIN feeds f ON f.id = a.feed_id
		WHERE a.id = ?
	`, id).Scan(
		&a.ID, &a.FeedID, &a.GUID, &a.Title, &a.URL, &a.Author,
		&a.Content, &a.Summary, &publishedAt, &a.FetchedAt, &a.Read, &readAt,
		&a.FeedTitle,
	)
	if err != nil {
		return nil, err
	}
	if publishedAt.Valid {
		a.PublishedAt = &publishedAt.Time
	}
	if readAt.Valid {
		a.ReadAt = &readAt.Time
	}
	return &a, nil
}

// MarkArticleRead marks an article as read.
func (db *DB) MarkArticleRead(id int64) error {
	_, err := db.conn.Exec(
		`UPDATE articles SET read = 1, read_at = CURRENT_TIMESTAMP WHERE id = ?`, id,
	)
	return err
}

// MarkAllRead marks all articles in a feed as read. If feedID is 0, marks all articles.
func (db *DB) MarkAllRead(feedID int64) error {
	if feedID > 0 {
		_, err := db.conn.Exec(
			`UPDATE articles SET read = 1, read_at = CURRENT_TIMESTAMP WHERE feed_id = ? AND read = 0`,
			feedID,
		)
		return err
	}
	_, err := db.conn.Exec(
		`UPDATE articles SET read = 1, read_at = CURRENT_TIMESTAMP WHERE read = 0`,
	)
	return err
}
