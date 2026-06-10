package parser

import (
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
	h := CreateHeaders()
	headerBuff := [12]byte(b[:12])
	err := h.Parse(headerBuff)
	if err != nil {
		return fmt.Errorf("failed to parse header prev: %s", err)
	}

	q := CreateQuestion()
	qBuff := b[12:]
	err = q.ParseQuestion(qBuff)
	if err != nil {
		return fmt.Errorf("failed to parse question prev: %s", err)
	}

	d.Header = *h
	d.Question = *q
	return nil
}
