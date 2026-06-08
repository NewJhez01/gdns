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

func Resolve(b []byte, c cache.Cache, bl blocklist.Blocklist) ([]byte, error) {
	dns := parser.CreateNewDnsStruct()
	err := dns.Parse(b)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	val, err := c.GetDomainNameFromCache(ctx, dns.Question.Qname)
	if errors.Is(err, cache.ErrEmpty) {
		return handleDns(dns.Question.Qname, bl, c)
	}
	if err != nil {
		return nil, err
	}
	if val.IsBlocked == true {
		// to do after the dns answer parser is built parse the reject into proper
		// resp and return that
		return []byte(REJECT), nil
	}
	return nil, nil
}

func handleDns(s string, b blocklist.Blocklist, c cache.Cache) ([]byte, error) {
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

	answer, err := forwardToUpstream()
	if err == nil {
		return nil, err
	}

	val := cache.Value{
		Answer:    answer,
		IsBlocked: false,
	}
	err = c.SetDomainName(ctx, s, val, 2*time.Minute)
	if err != nil {
		return nil, err
	}

	return answer, nil
}
