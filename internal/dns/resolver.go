package dns

import (
	"context"

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
		return handleDns(dns.Question.Qname, bl)
	}
	return val, nil
}

func handleDns(s string, b blocklist.Blocklist) (string, error) {
	isBlocked, err := b.IsBlocked(s)
	if err != nil {
		return "", err
	}

	if isBlocked == true {
		return REJECT, nil
	}

	// todo handle the fetching of the proper resolved address
	return "", nil
}
