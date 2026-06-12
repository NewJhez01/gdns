package unit

import (
	"testing"

	"github.com/gdns/internal/dns/parser"
	"github.com/stretchr/testify/assert"
)

func TestNxDomainParser(t *testing.T) {
	h := parser.Header{
		Aa:      true,
		Tc:      true,
		Rd:      true,
		Ra:      true,
		Z:       0,
		Id:      parser.SixteenBit{0x00, 0x01},
		QdCount: parser.SixteenBit{0x00, 0x01},
		AnCount: parser.SixteenBit{0x00, 0x01},
		NsCount: parser.SixteenBit{0x00, 0x01},
		ArCount: parser.SixteenBit{0x00, 0x01},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	q := parser.Question{
		Qname:  "a.example.com",
		Qtype:  parser.SixteenBit{0x00, 0x01},
		Qclass: parser.SixteenBit{0x00, 0x01},
	}
	d := parser.Dns{
		Header:   h,
		Question: q,
	}
	res, err := d.BuildNxDomainResp()
	expec := []byte{
		0x0, 0x1, 0x87, 0x83, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x61, 0x7, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x3, 0x63, 0x6f, 0x6d, 0x0, 0x0, 0x1, 0x0, 0x1,
	}
	assert.NoError(t, err)
	assert.Equal(t, expec, res)
}

func TestNxDomainParser_LabelTooLong(t *testing.T) {
	h := parser.Header{
		Aa:      true,
		Tc:      true,
		Rd:      true,
		Ra:      true,
		Z:       0,
		Id:      parser.SixteenBit{0x00, 0x01},
		QdCount: parser.SixteenBit{0x00, 0x01},
		AnCount: parser.SixteenBit{0x00, 0x01},
		NsCount: parser.SixteenBit{0x00, 0x01},
		ArCount: parser.SixteenBit{0x00, 0x01},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	q := parser.Question{
		Qname:  "a.exampleasjkdhfjkasdhfasdjkfhjkasdhfjkasdhjkfhasjkfhjkasdhfjkasdh.com",
		Qtype:  parser.SixteenBit{0x00, 0x01},
		Qclass: parser.SixteenBit{0x00, 0x01},
	}
	d := parser.Dns{
		Header:   h,
		Question: q,
	}
	_, err := d.BuildNxDomainResp()
	assert.EqualError(t, err, "label is too big")
}
