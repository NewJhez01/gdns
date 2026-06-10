package parser

import (
	"encoding/binary"
	"errors"
	"strings"
)

func (d *Dns) marshall() ([]byte, error) {
	b := make([]byte, 12)
	var misc uint16
	misc |= (uint16(d.Header.Qr) & 0x1) << 15
	misc |= (uint16(d.Header.OpCode) & 0xF) << 11
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
	parts := strings.SplitSeq(d.Question.Qname, ".")
	for v := range parts {
		if len(v) > 63 {
			return nil, errors.New("label is too big")
		}
		b = append(b, byte(len(v)))
		b = append(b, []byte(v)...)
	}
	b = append(b, 0x00)

	b = append(b, d.Question.Qtype[0])
	b = append(b, d.Question.Qtype[1])
	b = append(b, d.Question.Qclass[0])
	b = append(b, d.Question.Qclass[1])
	return b, nil
}

func (d *Dns) BuildNxDomainResp() ([]byte, error) {
	d.Header.Qr = Response
	d.Header.Rcode = NameErr
	d.Header.QdCount = SixteenBit{0x00, 0x01}
	d.Header.AnCount = SixteenBit{0x00, 0x00}
	d.Header.NsCount = SixteenBit{0x00, 0x00}
	d.Header.ArCount = SixteenBit{0x00, 0x00}
	return d.marshall()
}
