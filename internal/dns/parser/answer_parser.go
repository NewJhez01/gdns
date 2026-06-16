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
	nameLen, err := CountNameLen(b[offset:], 0)
	if err != nil {
		return err
	}
	if err := a.ParseName(b, 12+questionLen, 0); err != nil {
		return err
	}
	if offset+nameLen >= len(b) {
		return errors.New("malformed request")
	}
	if err := a.parseAfterAnswer(b[offset+nameLen:]); err != nil {
		return err
	}
	return nil
}

func CountNameLen(b []byte, consumed int) (int, error) {
	if len(b) == 0 {
		return consumed, nil
	}
	if b[0]&0xC0 == 0xC0 {
		return consumed + 2, nil
	}

	label := int(b[0])
	if label == 0 {
		return consumed + 1, nil
	}
	if len(b) <= label+1 {
		return 0, errors.New("malformed request")
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

func (a *Answer) ParseName(b []byte, cursor, depth int) error {
	if depth > 10 {
		return errors.New("endless pointer caught")
	}
	if cursor >= len(b) {
		return errors.New("malformed request")
	}
	if b[cursor] == 0 {
		return nil
	}
	// pointer found
	if b[cursor]&0xC0 == 0xC0 {
		if cursor+1 >= len(b) {
			return errors.New("malformed request")
		}
		restBits := binary.BigEndian.Uint16(b[cursor:])
		if err := a.ParseName(b, int(restBits&0x3FFF), depth+1); err != nil {
			return err
		}
		return nil
	}
	n, err := a.appendDomain(b[cursor:])
	if err != nil {
		return err
	}
	if err := a.ParseName(b, cursor+n, depth); err != nil {
		return err
	}
	return nil
}

func (a *Answer) appendDomain(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, errors.New("malformed request")
	}
	label := int(b[0])
	if len(b) <= label+1 {
		return 0, errors.New("malformed request")
	}
	tempstr := string(b[1 : label+1])
	a.Name = append(a.Name, tempstr)
	return len(tempstr) + 1, nil
}
