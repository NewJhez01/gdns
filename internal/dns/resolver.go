package dns

import (
	"context"
	"errors"
	"log"
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
	log.Print(dns.Question.Qname)
	ctx := context.Background()
	val, err := c.GetDomainNameFromCache(ctx, dns.Question.Qname)
	if errors.Is(err, cache.ErrEmpty) {
		return handleDns(b, dns.Question.Qname, bl, c, dns)
	}
	if err != nil {
		return nil, err
	}
	if val.IsBlocked {
		return dns.BuildNxDomainResp()
	}
	return parser.BuildResponse(*dns, val.Rdata, val.TTL)
}

func handleDns(buff []byte, s string, bl blocklist.Blocklist, c cache.Cache, d *parser.Dns) ([]byte, error) {
	ctx := context.Background()
	ctxWIthTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	isBlocked, err := bl.IsBlocked(s, ctxWIthTimeout)
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
		return d.BuildNxDomainResp()
	}

	answer, err := forwardToUpstream(buff, d.Question.Len)
	if err != nil {
		return nil, err
	}

	val := cache.Value{
		Rdata:     answer.Rdata,
		IsBlocked: false,
		TTL:       answer.Ttl,
	}
	if answer.Ttl > 0 {
		if err := c.SetDomainName(ctx, s, val, time.Duration(answer.Ttl)*time.Second); err != nil {
			return nil, err
		}
	}

	return parser.BuildResponse(*d, answer.Rdata, answer.Ttl)
}
