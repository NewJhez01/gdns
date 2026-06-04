package parser

import "fmt"

type Dns struct {
	Header   Header
	Question Question
}

func CreateNewDnsStruct() *Dns {
	return &Dns{}
}

func (d *Dns) Parse(b [512]byte) error {
	h := createHeaders()
	headerBuff := [12]byte(b[:12])
	err := h.Parse(headerBuff)
	if err != nil {
		return fmt.Errorf("failed to parse header prev: %s", err)
	}

	q := createQuestion()
	qBuff := b[12:]
	err = q.parseQuestion(qBuff)
	if err != nil {
		return fmt.Errorf("failed to parse question prev: %s", err)
	}

	d.Header = *h
	d.Question = *q
	return nil
}
