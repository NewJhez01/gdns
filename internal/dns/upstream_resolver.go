package dns

import (
	"fmt"
	"net"
)

func forwardToUpstream() ([]byte, error) {
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return nil, fmt.Errorf("failed to open conn prev: %s", err)
	}
	defer conn.Close()

	b := make([]byte, 512)
	n, err := conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read from conn prev: %s", err)
	}
	return b[:n], nil
}
