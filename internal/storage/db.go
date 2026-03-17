package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection.
type DB struct {
	conn *sql.DB
}

// Open creates or opens the SQLite database at the given path.
func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(on)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return db, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying *sql.DB for direct queries.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS documents (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		path        TEXT NOT NULL UNIQUE,
		format      TEXT NOT NULL CHECK(format IN ('pdf', 'epub')),
		author      TEXT DEFAULT '',
		total_pages INTEGER DEFAULT 0,
		file_size   INTEGER DEFAULT 0,
		added_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
		opened_at   DATETIME
	);

	CREATE TABLE IF NOT EXISTS reading_progress (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		document_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
		page        INTEGER DEFAULT 0,
		scroll_pct  REAL DEFAULT 0.0,
		updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(document_id)
	);

	CREATE TABLE IF NOT EXISTS highlights (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		document_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
		page        INTEGER NOT NULL,
		text        TEXT NOT NULL,
		note        TEXT DEFAULT '',
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS settings (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS feeds (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		title        TEXT NOT NULL,
		url          TEXT NOT NULL UNIQUE,
		site_url     TEXT DEFAULT '',
		description  TEXT DEFAULT '',
		image_url    TEXT DEFAULT '',
		last_fetched DATETIME,
		added_at     DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS articles (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_id      INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
		guid         TEXT NOT NULL,
		title        TEXT NOT NULL,
		url          TEXT DEFAULT '',
		author       TEXT DEFAULT '',
		content      TEXT DEFAULT '',
		summary      TEXT DEFAULT '',
		published_at DATETIME,
		fetched_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
		read         INTEGER DEFAULT 0,
		read_at      DATETIME,
		UNIQUE(feed_id, guid)
	);

	CREATE INDEX IF NOT EXISTS idx_articles_feed_id ON articles(feed_id);
	CREATE INDEX IF NOT EXISTS idx_articles_read ON articles(read);
	CREATE INDEX IF NOT EXISTS idx_articles_published ON articles(published_at DESC);
	`
	_, err := db.conn.Exec(schema)
	return err
}
