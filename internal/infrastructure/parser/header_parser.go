package parser

import (
	"encoding/binary"
	"fmt"
	"io"
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

func CreateHeader() *Header {
	// quote from rfc 1035
	// Reserved for future use.  Must be zero in all queries and responses.
	return &Header{
		Z: 0,
	}
}

func (h *Header) Parse(r io.Reader) error {
	buff := make([]byte, 12)
	_, err := io.ReadFull(r, buff)
	if err != nil {
		return fmt.Errorf("failed to read into buffer prev: %s", err)
	}
	h.Id = [2]byte(buff[:2])
	h.QdCount = [2]byte(buff[4:6])
	h.AnCount = [2]byte(buff[6:8])
	h.NsCount = [2]byte(buff[8:10])
	h.ArCount = [2]byte(buff[10:])

	miscBuf := buff[2:4]
	err = h.parseMisc(miscBuf)
	if err != nil {
		return fmt.Errorf("failed to parse the second row of 16 bits from dns prev: %s", err)
	}

	return nil
}

func (h *Header) parseMisc(b []byte) error {
	bits := binary.BigEndian.Uint16(b)
	h.Qr = QR((bits >> 15) & 0x1)
	h.Aa = (bits>>10)&0x1 == 0x1
	h.Tc = (bits>>9)&0x1 == 0x1
	h.Rd = (bits>>8)&0x1 == 0x1
	h.Ra = (bits>>7)&0x1 == 0x1

	return nil
}
