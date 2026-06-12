package integration

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/gdns/cmd/server"
	"github.com/gdns/internal/blocklist"
	"github.com/gdns/internal/cache"
	"github.com/redis/go-redis/v9"
)

func setupResolver(t *testing.T) (cleanup func()) {
	// use in memory sqlite instead of embedded db
	db, err := blocklist.CreateNewDbConn(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	if _, err := db.Db.Exec("INSERT INTO blocked_domains (domain) VALUES ('doubleclick.net')"); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	r := cache.CreateNewRedisClient(redisClient)

	go server.Serve(*r, *db)

	time.Sleep(100 * time.Millisecond)

	return func() {
		db.Db.Close()
		redisClient.Close()
	}
}

func TestResolver_NonBlocked(t *testing.T) {
	cleanup := setupResolver(t)
	defer cleanup()

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp", "127.0.0.1:5555")
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ips, err := r.LookupHost(ctx, "google.com")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}

	if len(ips) == 0 {
		t.Fatal("expected IPs, got none")
	}

	t.Logf("google.com resolved to: %v", ips)
}

func TestResolver_Blocked(t *testing.T) {
	cleanup := setupResolver(t)
	defer cleanup()

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp", "127.0.0.1:5555")
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.LookupHost(ctx, "doubleclick.net")
	if err == nil {
		t.Fatal("expected error for blocked domain, got nil")
	}

	// DNS NXDOMAIN comes back as "no such host"
	if dnsErr, ok := err.(*net.DNSError); ok {
		if !dnsErr.IsNotFound {
			t.Fatalf("expected NXDOMAIN, got: %v", err)
		}
		t.Logf("doubleclick.net blocked: %v", err)
	} else {
		t.Fatalf("unexpected error type: %T: %v", err, err)
	}
}
