package parser

import (
	"encoding/binary"
	"errors"
)

func BuildResponse(d Dns, rData []byte, ttl uint32) ([]byte, error) {
	d.Header.Qr = Response
	d.Header.Rcode = NoErrorCon
	d.Header.QdCount = SixteenBit{0x00, 0x01}
	d.Header.AnCount = SixteenBit{0x00, 0x01}
	d.Header.NsCount = SixteenBit{0x00, 0x00}
	d.Header.ArCount = SixteenBit{0x00, 0x00}
	question, err := d.marshall()
	if err != nil {
		return nil, errors.New("failed to marshall response")
	}
	ans := createDefaultAnswer(rData, ttl)
	buff := []byte{}
	buff = append(buff, question...)
	buff = append(buff, ans...)

	return buff, nil
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

func createDefaultAnswer(rData []byte, ttl uint32) []byte {
	buff := []byte{
		0xC0, 0x0C, // NAME pointer
		0x00, 0x01, // TYPE A
		0x00, 0x01, // CLASS IN
	}
	ttlBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlBytes, ttl)
	buff = append(buff, ttlBytes...)
	buff = append(buff, 0x00, 0x04) // RDLENGTH
	buff = append(buff, rData...)   // IP
	return buff
}
