package storage

// GetSetting retrieves a setting value by key. Returns empty string if not found.
func (db *DB) GetSetting(key string) string {
	var value string
	err := db.conn.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

// SetSetting stores a setting key-value pair.
func (db *DB) SetSetting(key, value string) error {
	_, err := db.conn.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}
