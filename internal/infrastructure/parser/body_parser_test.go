package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateQuestionSection1(t *testing.T) {
	// a.example.com
	b := []byte{
		0x01, 'a', 0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x03,
		'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01,
	}
	q := createQuestion()
	err := q.parseQuestion(b)
	assert.Nil(t, err)
	assert.Equal(t, "a.example.com", q.Qname)
	assert.Equal(t, [2]byte{0x00, 0x01}, q.Qtype)
	assert.Equal(t, [2]byte{0x00, 0x01}, q.Qclass)
}
