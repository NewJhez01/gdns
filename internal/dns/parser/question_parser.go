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
	Qtype  [2]byte
	Qclass [2]byte
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
		q.Qtype = [2]byte(b[:2])
		return 2, QCLASS, nil
	case QCLASS:
		q.Qclass = [2]byte(b[:2])
		return 2, DONE, nil
	}
	return 0, DONE, errors.New("invalid state")
}
