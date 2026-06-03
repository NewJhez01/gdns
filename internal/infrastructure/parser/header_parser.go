package parser

import (
	"encoding/binary"
	"fmt"
)

type (
	QR     uint8
	OPCODE uint8
	RCODE  uint8
)

const (
	query QR = iota
	response
)

const (
	standardQuery OPCODE = iota
	inverseQuery
	serverStatus
)

const (
	noErrorCon RCODE = iota
	formatErr
	serverFailure
	nameErr
	notImplemented
	refused
)

// see rfc 1035 4.1.1 for exact header definition
// https://datatracker.ietf.org/doc/html/rfc1035#autoid-41
type Header struct {
	Aa      bool
	Tc      bool
	Rd      bool
	Ra      bool
	Z       uint8
	Id      [2]byte
	QdCount [2]byte
	AnCount [2]byte
	NsCount [2]byte
	ArCount [2]byte
	OpCode  OPCODE
	Qr      QR
	Rcode   RCODE
}

func CreateHeaders() *Header {
	return &Header{}
}

func (h *Header) Parse(buff [12]byte) error {
	h.Id = [2]byte(buff[:2])
	h.QdCount = [2]byte(buff[4:6])
	h.AnCount = [2]byte(buff[6:8])
	h.NsCount = [2]byte(buff[8:10])
	h.ArCount = [2]byte(buff[10:])

	miscBuf := buff[2:4]
	err := h.parseMisc(miscBuf)
	if err != nil {
		return fmt.Errorf("failed to parse the second row of 16 bits from dns prev: %s", err)
	}

	return nil
}

func (h *Header) parseMisc(b []byte) error {
	bits := binary.BigEndian.Uint16(b)
	h.Qr = QR((bits >> 15) & 0x1)
	h.Aa = (bits>>10)&0x1 == 1
	h.Tc = (bits>>9)&0x1 == 1
	h.Rd = (bits>>8)&0x1 == 1
	h.Ra = (bits>>7)&0x1 == 1
	h.OpCode = OPCODE((bits >> 11) & 0xF)
	h.Rcode = RCODE(bits & 0xF)
	if z := ((bits >> 4) & 0x7); z != 0 {
		return fmt.Errorf("Z must be 0 and is %d", z)
	}

	return nil
}
