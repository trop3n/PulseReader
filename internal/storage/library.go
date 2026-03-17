package storage

import (
	"database/sql"
	"time"
)

// Document represents a PDF or EPUB in the library.
type Document struct {
	ID         int64
	Title      string
	Path       string
	Format     string
	Author     string
	TotalPages int
	FileSize   int64
	AddedAt    time.Time
	OpenedAt   *time.Time

	// Joined from reading_progress
	CurrentPage int
	ScrollPct   float64
}

// AddDocument adds a document to the library.
func (db *DB) AddDocument(doc *Document) error {
	res, err := db.conn.Exec(
		`INSERT INTO documents (title, path, format, author, total_pages, file_size)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET title=excluded.title, total_pages=excluded.total_pages`,
		doc.Title, doc.Path, doc.Format, doc.Author, doc.TotalPages, doc.FileSize,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	if id > 0 {
		doc.ID = id
	}
	return nil
}

// ListDocuments returns all documents ordered by most recently opened.
func (db *DB) ListDocuments() ([]Document, error) {
	rows, err := db.conn.Query(`
		SELECT d.id, d.title, d.path, d.format, d.author, d.total_pages, d.file_size,
		       d.added_at, d.opened_at,
		       COALESCE(rp.page, 0), COALESCE(rp.scroll_pct, 0.0)
		FROM documents d
		LEFT JOIN reading_progress rp ON rp.document_id = d.id
		ORDER BY COALESCE(d.opened_at, d.added_at) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var d Document
		var openedAt sql.NullTime
		err := rows.Scan(
			&d.ID, &d.Title, &d.Path, &d.Format, &d.Author, &d.TotalPages, &d.FileSize,
			&d.AddedAt, &openedAt,
			&d.CurrentPage, &d.ScrollPct,
		)
		if err != nil {
			return nil, err
		}
		if openedAt.Valid {
			d.OpenedAt = &openedAt.Time
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

// GetDocumentByPath retrieves a document by its file path.
func (db *DB) GetDocumentByPath(path string) (*Document, error) {
	var d Document
	var openedAt sql.NullTime
	err := db.conn.QueryRow(`
		SELECT d.id, d.title, d.path, d.format, d.author, d.total_pages, d.file_size,
		       d.added_at, d.opened_at,
		       COALESCE(rp.page, 0), COALESCE(rp.scroll_pct, 0.0)
		FROM documents d
		LEFT JOIN reading_progress rp ON rp.document_id = d.id
		WHERE d.path = ?
	`, path).Scan(
		&d.ID, &d.Title, &d.Path, &d.Format, &d.Author, &d.TotalPages, &d.FileSize,
		&d.AddedAt, &openedAt,
		&d.CurrentPage, &d.ScrollPct,
	)
	if err != nil {
		return nil, err
	}
	if openedAt.Valid {
		d.OpenedAt = &openedAt.Time
	}
	return &d, nil
}

// TouchDocument updates the opened_at timestamp.
func (db *DB) TouchDocument(id int64) error {
	_, err := db.conn.Exec(`UPDATE documents SET opened_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

// RemoveDocument deletes a document from the library.
func (db *DB) RemoveDocument(id int64) error {
	_, err := db.conn.Exec(`DELETE FROM documents WHERE id = ?`, id)
	return err
}
