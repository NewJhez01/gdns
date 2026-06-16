package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/gdns/cmd/server"
	"github.com/gdns/internal/blocklist"
	"github.com/gdns/internal/cache"
	"github.com/redis/go-redis/v9"
)

func setupResolver(t *testing.T) (resolver *net.Resolver, cleanup func()) {
	os.Setenv("UDP_PORT", "127.0.0.1:0")

	redisAddr := os.Getenv("REDIS_URL")
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr, DB: 1})

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
		t.Fatalf("failed to insert blocked domain: %v", err)
	}

	srv := server.NewServer(*cache.CreateNewRedisClient(redisClient), *db)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	udpAddr, ok := srv.Conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		t.Fatalf("expected UDPAddr, got %T", srv.Conn.LocalAddr())
	}
	port := udpAddr.Port

	resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
		},
	}

	cleanup = func() {
		srv.Stop()
		redisClient.Close()
		db.Db.Close()
		redisClient.FlushDB(ctx)
	}

	return resolver, cleanup
}

func TestResolver(t *testing.T) {
	r, cleanup := setupResolver(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test non-blocked
	ips, err := r.LookupHost(ctx, "google.com")
	if err != nil {
		t.Fatalf("non-blocked lookup failed: %v", err)
	}
	if len(ips) == 0 {
		t.Fatal("expected IPs for google.com, got none")
	}
	t.Logf("google.com -> %v", ips)

	// Test blocked
	_, err = r.LookupHost(ctx, "doubleclick.net")
	if err == nil {
		t.Fatal("expected error for blocked domain, got nil")
	}
	if dnsErr, ok := err.(*net.DNSError); !ok || !dnsErr.IsNotFound {
		t.Fatalf("expected NXDOMAIN for doubleclick.net, got: %v", err)
	}
	t.Logf("doubleclick.net -> blocked (NXDOMAIN)")
}
