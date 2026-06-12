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
	Db *sql.DB
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
	err := r.Db.QueryRowContext(
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
	_, err := r.Db.ExecContext(ctx,
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
	defer closeConn(f)
	sc := bufio.NewScanner(f)
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT OR IGNORE INTO blocked_domains(domain) VALUES(?)")
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatalf("failed to rollback transcation prev: %s", err)
		}
		return err
	}
	defer closeStmt(stmt)

	count := 0
	for sc.Scan() {
		domain := strings.TrimSpace(sc.Text())
		if domain == "" {
			continue
		}
		_, err := stmt.ExecContext(ctx, domain)
		if err != nil {
			log.Printf("failed to insert %s: %v", domain, err)
			continue
		}
		count++
	}
	log.Printf("inserted %d domains", count)

	if err := sc.Err(); err != nil {
		log.Printf("scanner error: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Fatalf("failed to rollback transcation prev: %s", err)
		}
		return err
	}
	return tx.Commit()
}

func closeConn(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatalf("failed to close db prev: %s", err)
	}
}

func closeStmt(s *sql.Stmt) {
	if err := s.Close(); err != nil {
		log.Fatalf("failed to close statement prev %s", err)
	}
}
