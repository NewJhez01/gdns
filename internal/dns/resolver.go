package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/gdns/internal/blocklist"
	"github.com/gdns/internal/cache"
	"github.com/gdns/internal/dns/parser"
)

const REJECT = "NXDOMAIN"

func Resolve(b []byte, c cache.Cache, bl blocklist.Blocklist) (string, error) {
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
		return REJECT, nil
	}
	if val == "" {
		return handleDns(dns.Question.Qname, bl, c)
	}
	return val, nil
}

func handleDns(s string, b blocklist.Blocklist, c cache.Cache) (string, error) {
	ctx := context.Background()
	ctxWIthTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	isBlocked, err := b.IsBlocked(s, ctxWIthTimeout)
	if err != nil {
		return "", err
	}

	if isBlocked {
		c.SetDomainName(ctx, s, fmt.Sprintf("%t", isBlocked), 15*time.Minute)
		return REJECT, nil
	}

	// todo handle the fetching of the proper resolved address

	c.SetDomainName(ctx, s, fmt.Sprintf("%t", isBlocked), 2*time.Minute)

	return "", nil
}
