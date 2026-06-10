package parser

import (
	"encoding/binary"
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

func (d *Dns) Marshall() []byte {
	b := [12]byte{}
	var misc uint16
	misc |= uint16((d.Header.Qr)&0x1) << 15
	misc |= uint16((d.Header.OpCode)&0x7) << 11
	if d.Header.Aa {
		misc |= 1 << 10
	}
	if d.Header.Tc {
		misc |= 1 << 9
	}
	if d.Header.Rd {
		misc |= 1 << 8
	}
	if d.Header.Ra {
		misc |= 1 << 7
	}
	misc |= uint16((d.Header.Z)&0x7) << 4
	misc |= uint16(d.Header.Rcode) & 0xF
	copy(b[0:2], d.Header.Id[:])
	binary.BigEndian.PutUint16(b[2:4], misc)
	copy(b[4:6], d.Header.QdCount[:])
	copy(b[6:8], d.Header.AnCount[:])
	copy(b[8:10], d.Header.NsCount[:])
	copy(b[10:12], d.Header.ArCount[:])
	return b[:]
}

func (d *Dns) BuildNxDomainResp() []byte {
	d.Header.Qr = Response
	d.Header.Rcode = NameErr
	d.Header.QdCount = SixteenBit{1}
	d.Header.AnCount = SixteenBit{0}
	d.Header.NsCount = SixteenBit{0}
	d.Header.ArCount = SixteenBit{0}
	return d.Marshall()
}
