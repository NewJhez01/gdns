package dns

import (
	"context"
	"errors"
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
	if errors.Is(err, cache.ErrEmpty) {
		return handleDns(dns.Question.Qname, bl, c)
	}
	if err != nil {
		return "", err
	}
	if val.IsBlocked == true {
		return REJECT, nil
	}
	return "", nil
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
		val := cache.Value{
			Ip:        REJECT,
			IsBlocked: true,
		}
		err := c.SetDomainName(ctx, s, val, 15*time.Minute)
		if err != nil {
			return "", err
		}
		return REJECT, nil
	}

	// todo handle the fetching of the proper resolved address
	// this is also where the ip is then getting fetched to set in cache
	val := cache.Value{
		Ip:        "",
		IsBlocked: false,
	}
	err = c.SetDomainName(ctx, s, val, 2*time.Minute)
	if err != nil {
		return "", err
	}

	return "", nil
}
