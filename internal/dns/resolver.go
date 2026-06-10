package dns

import (
	"context"
	"errors"
	"time"

	"github.com/gdns/internal/blocklist"
	"github.com/gdns/internal/cache"
	"github.com/gdns/internal/dns/parser"
)

func Resolve(b []byte, c cache.Cache, bl blocklist.Blocklist) ([]byte, error) {
	dns := parser.CreateNewDnsStruct()
	if err := dns.Parse(b); err != nil {
		return nil, err
	}
	ctx := context.Background()
	val, err := c.GetDomainNameFromCache(ctx, dns.Question.Qname)
	if errors.Is(err, cache.ErrEmpty) {
		return handleDns(b, dns.Question.Qname, bl, c)
	}
	if err != nil {
		return nil, err
	}
	if val.IsBlocked {
		return dns.BuildNxDomainResp(), nil
	}
	return nil, nil
}

func handleDns(buff []byte, s string, b blocklist.Blocklist, c cache.Cache) ([]byte, error) {
	ctx := context.Background()
	ctxWIthTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	isBlocked, err := b.IsBlocked(s, ctxWIthTimeout)
	if err != nil {
		return nil, err
	}

	if isBlocked {
		val := cache.Value{
			IsBlocked: true,
		}
		err := c.SetDomainName(ctx, s, val, 15*time.Minute)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	answer, err := forwardToUpstream(buff)
	if err != nil {
		return nil, err
	}

	val := cache.Value{
		Answer:    answer,
		IsBlocked: false,
	}
	if err := c.SetDomainName(ctx, s, val, 2*time.Minute); err != nil {
		return nil, err
	}

	return answer, nil
}
