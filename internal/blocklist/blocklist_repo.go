package blocklist

import (
	"database/sql"
	"fmt"
)

type SqliteClient struct {
	db *sql.DB
}

func CreateNewDbConn(path string) (*SqliteClient, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database prev: %s", err)
	}
	return &SqliteClient{db}, nil
}

func (r *SqliteClient) IsBlocked(key string) (bool, error) {
	isBlocked := false
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM blocked_domains WHERE domain = ?)",
		key,
	).Scan(&isBlocked)
	if err != nil {
		return isBlocked, fmt.Errorf("failed to read db prev: %s", err)
	}
	return isBlocked, nil
}

func (r *SqliteClient) Migrate() error {
	_, err := r.db.Exec(`
        CREATE TABLE IF NOT EXISTS blocked_domains (
            domain TEXT PRIMARY KEY
        )
    `)
	if err != nil {
		return err
	}
	return nil
}
