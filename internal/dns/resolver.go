package dns

import (
	"github.com/gdns/internal/dns/parser"
)

func Resolve(b []byte) error {
	dnsStruct := parser.CreateNewDnsStruct()
	err := dnsStruct.Parse((b))
	if err != nil {
		return err
	}
	return nil
}
