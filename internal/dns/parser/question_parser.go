package parser

import (
	"errors"
)

type state int

const (
	QNAME state = iota
	QTYPE
	QCLASS
	DONE
)

type Question struct {
	Qname  string
	Qtype  SixteenBit
	Qclass SixteenBit
}

func CreateQuestion() *Question {
	return &Question{}
}

func (q *Question) ParseQuestion(b []byte) error {
	state := QNAME
	bytesRead := 0
	for {
		if state == DONE {
			break
		}
		n, s, err := q.parseStruct(b[bytesRead:], state)
		if err != nil {
			return err
		}
		state = s
		bytesRead += n
	}
	return nil
}

func (q *Question) parseStruct(b []byte, s state) (int, state, error) {
	switch s {
	case QNAME:
		if b[0] == 0 {
			return 1, QTYPE, nil
		}
		if q.Qname != "" {
			q.Qname += "."
		}
		q.Qname += string(b[1 : 1+b[0]])

		return 1 + int(b[0]), QNAME, nil
	case QTYPE:
		if len(b) < 2 {
			return 0, 0, errors.New("malformed dns package")
		}
		q.Qtype = SixteenBit((b[:2]))
		return 2, QCLASS, nil
	case QCLASS:
		if len(b) < 2 {
			return 0, 0, errors.New("malformed dns package")
		}
		q.Qclass = SixteenBit(b[:2])
		return 2, DONE, nil
	}
	return 0, DONE, errors.New("invalid state")
}
