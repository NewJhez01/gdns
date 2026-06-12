package parser

import (
	"errors"
	"fmt"
)

type Dns struct {
	Header   Header
	Question Question
}

func CreateNewDnsStruct() *Dns {
	return &Dns{}
}

type SixteenBit [2]byte

func (d *Dns) Parse(b []byte) error {
	if len(b) < 12 {
		return errors.New("header is too short")
	}
	h := CreateHeaders()
	headerBuff := [12]byte(b[:12])
	h.Parse(headerBuff)

	q := CreateQuestion()
	qBuff := b[12:]
	err := q.ParseQuestion(qBuff)
	if err != nil {
		return fmt.Errorf("failed to parse question prev: %s", err)
	}

	d.Header = *h
	d.Question = *q
	return nil
}
