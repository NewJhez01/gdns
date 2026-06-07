package blocklist

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SqliteClient struct {
	db *sql.DB
}

func CreateNewDbConn(path string) (*SqliteClient, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database prev: %s", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("the db is unavailiblem, prev: %s", err)
	}
	return &SqliteClient{db}, nil
}

func (r *SqliteClient) IsBlocked(key string, ctx context.Context) (bool, error) {
	isBlocked := false
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM blocked_domains WHERE domain = ?)",
		key,
	).Scan(&isBlocked)
	if err != nil {
		return isBlocked, fmt.Errorf("failed to read db prev: %s", err)
	}
	return isBlocked, nil
}

func (r *SqliteClient) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		`
        CREATE TABLE IF NOT EXISTS blocked_domains (
            domain TEXT PRIMARY KEY
        )
    `)
	if err != nil {
		return err
	}
	return nil
}
