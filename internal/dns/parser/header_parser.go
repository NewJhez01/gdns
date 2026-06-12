package parser

import (
	"encoding/binary"
	"log"
)

type (
	QR     uint8
	OPCODE uint8
	RCODE  uint8
)

const (
	Query QR = iota
	Response
)

const (
	StandardQuery OPCODE = iota
	InverseQuery
	ServerStatus
)

const (
	NoErrorCon RCODE = iota
	FormatErr
	ServerFailure
	NameErr
	NotImplemented
	Refused
)

// see rfc 1035 4.1.1 for exact header definition
// https://datatracker.ietf.org/doc/html/rfc1035#autoid-41
type Header struct {
	Aa      bool
	Tc      bool
	Rd      bool
	Ra      bool
	Z       uint8
	Id      SixteenBit
	QdCount SixteenBit
	AnCount SixteenBit
	NsCount SixteenBit
	ArCount SixteenBit
	OpCode  OPCODE
	Qr      QR
	Rcode   RCODE
}

func CreateHeaders() *Header {
	return &Header{}
}

func (h *Header) Parse(buff [12]byte) {
	h.Id = SixteenBit(buff[:2])
	h.QdCount = SixteenBit(buff[4:6])
	h.AnCount = SixteenBit(buff[6:8])
	h.NsCount = SixteenBit(buff[8:10])
	h.ArCount = SixteenBit(buff[10:])

	miscBuf := buff[2:4]
	h.parseMisc(miscBuf)
}

func (h *Header) parseMisc(b []byte) {
	bits := binary.BigEndian.Uint16(b)
	h.Qr = QR((bits >> 15) & 0x1)
	h.OpCode = OPCODE((bits >> 11) & 0xF)
	h.Aa = (bits>>10)&0x1 == 1
	h.Tc = (bits>>9)&0x1 == 1
	h.Rd = (bits>>8)&0x1 == 1
	h.Ra = (bits>>7)&0x1 == 1
	if z := ((bits >> 4) & 0x7); z != 0 {
		log.Printf("Z should be 0 and is %d", z)
	}
	h.Rcode = RCODE(bits & 0xF)
}
