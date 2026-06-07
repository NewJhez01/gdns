package dns

import (
	"context"

	"github.com/gdns/internal/cache"
	"github.com/gdns/internal/dns/parser"
)

func Resolve(b []byte, c cache.Cache) (string, error) {
	dns := parser.CreateNewDnsStruct()
	err := dns.Parse(b)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	val, err := c.GetDomainNameFromCache(ctx, dns.Question.Qname)
	if err != nil {
		return "", err
	}
	if val == "blocked" {
		return "NXDOMAIN", nil
	}
	if val == "" {
		return handleDns(dns.Question.Qname), nil
	}
	return val, nil
}

func handleDns(s string) string {
	// todo implement logic for fetching the blocked status from the sql lite db
	return s
}
