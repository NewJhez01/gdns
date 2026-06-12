package parser

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Answer struct {
	Name     []string
	Type     SixteenBit
	Class    SixteenBit
	Ttl      uint32
	RdLength SixteenBit
	Rdata    []byte
}

func NewAnswer() *Answer {
	return &Answer{}
}

func (a *Answer) ParseAnswer(b []byte, questionLen int) error {
	offset := 12 + questionLen
	if len(b) < offset {
		return errors.New("answer section is too short")
	}
	nameLen := CountNameLen(b[offset:], 0)
	a.ParseName(b, 12+questionLen)
	if err := a.parseAfterAnswer(b[offset+nameLen:]); err != nil {
		return err
	}
	return nil
}
}

func CountNameLen(b []byte, consumed int) int {
	if len(b) == 0 {
		return consumed
	}
	if b[0]&0xC0 == 0xC0 {
		return consumed + 2
	}

	label := int(b[0])
	if label == 0 {
		return consumed + 1
	}
	return CountNameLen(b[label+1:], consumed+1+label)
}

func (a *Answer) parseAfterAnswer(b []byte) error {
	if len(b) < 10 {
		return errors.New("answer section is too short")
	}
	a.Type = SixteenBit{b[0], b[1]}
	a.Class = SixteenBit{b[2], b[3]}
	a.Ttl = binary.BigEndian.Uint32(b[4:8])
	a.RdLength = SixteenBit{b[8], b[9]}
	rdlength := binary.BigEndian.Uint16(a.RdLength[:])
	if len(b) < 10+int(rdlength) {
		return fmt.Errorf("rd length section is too short")
	}
	a.Rdata = b[10 : 10+rdlength]
	return nil
}

func (a *Answer) ParseName(b []byte, cursor int) {
	if b[cursor] == 0 {
		return
	}
	if b[cursor]&0xC0 == 0xC0 {
		// ptr found call recursive with offset of pointer
		restBits := binary.BigEndian.Uint16(b[cursor:])
		a.ParseName(b, int(restBits&0x3FFF))
		return
	}
	// parse label plus name
	n := a.appendDomain(b[cursor:])
	a.ParseName(b, cursor+n)
}

func (a *Answer) appendDomain(b []byte) int {
	label := int(b[0])
	tempstr := string(b[1 : label+1])
	a.Name = append(a.Name, tempstr)
	return len(tempstr) + 1
}
