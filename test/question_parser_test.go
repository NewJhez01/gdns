package test

import (
	"testing"

	"github.com/gdns/internal/dns/parser"
	"github.com/stretchr/testify/assert"
)

func TestCreateQuestionSection1(t *testing.T) {
	// a.example.com
	b := []byte{
		0x01, 'a', 0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x03,
		'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01,
	}
	q := parser.CreateQuestion()
	err := q.ParseQuestion(b)
	assert.NoError(t, err)
	assert.Equal(t, "a.example.com", q.Qname)
	assert.Equal(t, parser.SixteenBit{0x00, 0x01}, q.Qtype)
	assert.Equal(t, parser.SixteenBit{0x00, 0x01}, q.Qclass)
}
