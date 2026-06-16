package unit

import (
	"bytes"
	"net"
	"testing"

	"github.com/gdns/internal/dns/parser"
	"github.com/stretchr/testify/assert"
)

func TestParserFuzz(t *testing.T) {}

func TestCountNameLen(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int
	}{
		{
			name:     "pointer only",
			input:    []byte{0xC0, 0x0C},
			expected: 2,
		},
		{
			name:     "labels with null terminator",
			input:    []byte{3, 'w', 'w', 'w', 6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0},
			expected: 16,
		},
		{
			name:     "labels ending with pointer",
			input:    []byte{3, 'w', 'w', 'w', 0xC0, 0x0C},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.CountNameLen(tt.input, 0)
			if assert.NoError(t, err) == false {
				t.Errorf("unexpected error: %s", err)
			}
			if got != tt.expected {
				t.Errorf("CountNameLen() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestParseName(t *testing.T) {
	header := make([]byte, 12)
	questionName := []byte{3, 'w', 'w', 'w', 6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0}
	qtype := []byte{0x00, 0x01}
	qclass := []byte{0x00, 0x01}

	msg := append(header, questionName...)
	msg = append(msg, qtype...)
	msg = append(msg, qclass...)

	answer := []byte{0xC0, 0x0C, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x3C, 0x00, 0x04, 0xD8, 0x3A, 0xCF, 0x0E}
	msg = append(msg, answer...)

	t.Run("follows pointer to question name", func(t *testing.T) {
		a := parser.NewAnswer()
		err := a.ParseName(msg, 32, 0)
		if assert.NoError(t, err) == false {
			t.Errorf("unexpected error: %s", err)
		}

		want := []string{"www", "google", "com"}
		if len(a.Name) != len(want) {
			t.Fatalf("got %v, want %v", a.Name, want)
		}
		for i, label := range want {
			if a.Name[i] != label {
				t.Errorf("label %d: got %q, want %q", i, a.Name[i], label)
			}
		}
	})

	t.Run("reads direct name in answer section", func(t *testing.T) {
		directName := []byte{3, 'f', 'o', 'o', 0}
		rest := []byte{0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x3C, 0x00, 0x04, 0x08, 0x08, 0x08, 0x08}
		directMsg := append(header, questionName...)
		directMsg = append(directMsg, qtype...)
		directMsg = append(directMsg, qclass...)
		directMsg = append(directMsg, directName...)
		directMsg = append(directMsg, rest...)

		a := parser.NewAnswer()
		a.ParseName(directMsg, 32, 0)

		want := []string{"foo"}
		if len(a.Name) != len(want) || a.Name[0] != want[0] {
			t.Errorf("got %v, want %v", a.Name, want)
		}
	})
}

func TestParseAnswer(t *testing.T) {
	header := []byte{
		0x00, 0x01, // ID
		0x81, 0x80, // QR=1, RD=1, RA=1
		0x00, 0x01, // QDCOUNT
		0x00, 0x01, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
	}

	questionName := []byte{3, 'w', 'w', 'w', 6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0}
	qtype := []byte{0x00, 0x01}
	qclass := []byte{0x00, 0x01}
	questionLen := len(questionName) + len(qtype) + len(qclass) // 20

	answer := []byte{
		0xC0, 0x0C, // NAME pointer to offset 12
		0x00, 0x01, // TYPE A
		0x00, 0x01, // CLASS IN
		0x00, 0x00, 0x00, 0x3C, // TTL = 60
		0x00, 0x04, // RDLENGTH
		0xD8, 0x3A, 0xCF, 0x0E, // 216.58.207.14
	}

	msg := append(header, questionName...)
	msg = append(msg, qtype...)
	msg = append(msg, qclass...)
	msg = append(msg, answer...)

	t.Run("parses full A record answer", func(t *testing.T) {
		a := parser.NewAnswer()
		if err := a.ParseAnswer(msg, questionLen); err != nil {
			t.Fatalf("ParseAnswer() error = %v", err)
		}

		if a.Type != (parser.SixteenBit{0x00, 0x01}) {
			t.Errorf("Type = %v, want {0x00, 0x01}", a.Type)
		}
		if a.Class != (parser.SixteenBit{0x00, 0x01}) {
			t.Errorf("Class = %v, want {0x00, 0x01}", a.Class)
		}
		if a.Ttl != 60 {
			t.Errorf("Ttl = %d, want 60", a.Ttl)
		}
		if a.RdLength != (parser.SixteenBit{0x00, 0x04}) {
			t.Errorf("RdLength = %v, want {0x00, 0x04}", a.RdLength)
		}

		wantIP := net.IP{0xD8, 0x3A, 0xCF, 0x0E}
		if !bytes.Equal(a.Rdata, wantIP) {
			t.Errorf("Rdata = %v, want %v", a.Rdata, wantIP)
		}
	})

	t.Run("returns error on truncated answer", func(t *testing.T) {
		truncated := msg[:len(msg)-5] // chop off part of RDATA
		a := parser.NewAnswer()
		if err := a.ParseAnswer(truncated, questionLen); err == nil {
			t.Error("expected error for truncated answer, got nil")
		}
	})
}
