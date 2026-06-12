package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/gdns/internal/dns/parser"
)

func forwardToUpstream(query []byte, questionLen int) (*parser.Answer, error) {
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return nil, fmt.Errorf("failed to open conn prev: %s", err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(time.Second * 3)); err != nil {
		return nil, err
	}

	if _, err := conn.Write(query); err != nil {
		return nil, err
	}

	b := make([]byte, 512)
	n, err := conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read from conn prev: %s", err)
	}
	a := parser.NewAnswer()
	if err := a.ParseAnswer(b[:n], questionLen); err != nil {
		return a, err
	}
	return a, nil
}
