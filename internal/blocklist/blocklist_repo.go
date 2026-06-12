package blocklist

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
		return nil, fmt.Errorf("the db is unavailable, prev: %s", err)
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
        ) WITHOUT ROWID
    `)
	if err != nil {
		return err
	}

	f, err := os.Open(os.Getenv("SQLITE_BLOCKLIST"))
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	tx, err := r.db.Begin()
	stmt, err := tx.Prepare("INSERT OR IGNORE INTO blocked_domains(domain) VALUES(?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	count := 0
	for sc.Scan() {
		domain := strings.TrimSpace(sc.Text())
		if domain == "" {
			continue
		}
		_, err := stmt.Exec(domain)
		if err != nil {
			log.Printf("failed to insert %s: %v", domain, err)
			continue
		}
		count++
	}
	log.Printf("inserted %d domains", count)

	if err := sc.Err(); err != nil {
		log.Printf("scanner error: %v", err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
