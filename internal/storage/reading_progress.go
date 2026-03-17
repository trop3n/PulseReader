package storage

// SaveProgress saves reading position for a document.
func (db *DB) SaveProgress(documentID int64, page int, scrollPct float64) error {
	_, err := db.conn.Exec(`
		INSERT INTO reading_progress (document_id, page, scroll_pct, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(document_id) DO UPDATE SET
			page = excluded.page,
			scroll_pct = excluded.scroll_pct,
			updated_at = CURRENT_TIMESTAMP
	`, documentID, page, scrollPct)
	return err
}

// GetProgress returns the saved reading position for a document.
func (db *DB) GetProgress(documentID int64) (page int, scrollPct float64, err error) {
	err = db.conn.QueryRow(
		`SELECT page, scroll_pct FROM reading_progress WHERE document_id = ?`,
		documentID,
	).Scan(&page, &scrollPct)
	if err != nil {
		return 0, 0, nil // No progress saved yet
	}
	return page, scrollPct, nil
}
