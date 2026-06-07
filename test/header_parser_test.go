package test

import (
	"testing"

	"github.com/gdns/internal/dns/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser_StandardQuery(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      true,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0x00, 0x01},
		QdCount: parser.SixteenBit{0x00, 0x01},
		AnCount: parser.SixteenBit{0},
		NsCount: parser.SixteenBit{0},
		ArCount: parser.SixteenBit{0},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ResponseWithAllFlags(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0xAB, 0xCD, 0x87, 0x80, 0x00, 0x02, 0x00, 0x01, 0x00, 0x03, 0x00, 0x04}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      true,
		Tc:      true,
		Rd:      true,
		Ra:      true,
		Z:       0,
		Id:      parser.SixteenBit{0xAB, 0xCD},
		QdCount: parser.SixteenBit{0x00, 0x02},
		AnCount: parser.SixteenBit{0x00, 0x01},
		NsCount: parser.SixteenBit{0x00, 0x03},
		ArCount: parser.SixteenBit{0x00, 0x04},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Response,
		Rcode:   parser.NoErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_InverseQuery(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0x12, 0x34, 0x08, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0x12, 0x34},
		QdCount: parser.SixteenBit{0x00, 0x01},
		AnCount: parser.SixteenBit{0},
		NsCount: parser.SixteenBit{0},
		ArCount: parser.SixteenBit{0},
		OpCode:  parser.InverseQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ServerStatus(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0xFF, 0xFF, 0x90, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0xFF, 0xFF},
		QdCount: parser.SixteenBit{0},
		AnCount: parser.SixteenBit{0},
		NsCount: parser.SixteenBit{0},
		ArCount: parser.SixteenBit{0},
		OpCode:  parser.ServerStatus,
		Qr:      parser.Response,
		Rcode:   parser.ServerFailure,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_VariousRCodes(t *testing.T) {
	tests := []struct {
		name     string
		rcode    uint8
		expected parser.RCODE
	}{
		{"noError", 0, parser.NoErrorCon},
		{"formatError", 1, parser.FormatErr},
		{"ServerFailure", 2, parser.ServerFailure},
		{"NameError", 3, parser.NameErr},
		{"notImplemented", 4, parser.NotImplemented},
		{"refused", 5, parser.Refused},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := parser.CreateHeaders()
			secondRow := uint16(tt.rcode)
			dnsMock := [12]byte{
				0x00, 0x00,
				byte(secondRow >> 8), byte(secondRow & 0xFF),
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}
			err := h.Parse(dnsMock)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, h.Rcode)
		})
	}
}

func TestParser_MaxCounts(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0, 0},
		QdCount: parser.SixteenBit{0xFF, 0xFF},
		AnCount: parser.SixteenBit{0xFF, 0xFF},
		NsCount: parser.SixteenBit{0xFF, 0xFF},
		ArCount: parser.SixteenBit{0xFF, 0xFF},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ZeroBuffer(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0, 0},
		QdCount: parser.SixteenBit{0, 0},
		AnCount: parser.SixteenBit{0, 0},
		NsCount: parser.SixteenBit{0, 0},
		ArCount: parser.SixteenBit{0, 0},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ZNonZeroError(t *testing.T) {
	tests := []struct {
		name string
		val  byte
	}{
		{"z1", 0x10},
		{"z2", 0x20},
		{"z4", 0x40},
		{"z7", 0x70},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := parser.CreateHeaders()
			dnsMock := [12]byte{0x00, 0x00, 0x00, tt.val, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			err := h.Parse(dnsMock)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Z must be 0")
		})
	}
}

func TestParser_RealWorldGoogleDNS(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0x12, 0x34, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      true,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0x12, 0x34},
		QdCount: parser.SixteenBit{0x00, 0x01},
		AnCount: parser.SixteenBit{0},
		NsCount: parser.SixteenBit{0},
		ArCount: parser.SixteenBit{0},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Query,
		Rcode:   parser.NoErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ResponseWithNXDOMAIN(t *testing.T) {
	h := parser.CreateHeaders()
	dnsMock := [12]byte{0x00, 0x01, 0x80, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &parser.Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      parser.SixteenBit{0x00, 0x01},
		QdCount: parser.SixteenBit{0x00, 0x01},
		AnCount: parser.SixteenBit{0},
		NsCount: parser.SixteenBit{0x00, 0x01},
		ArCount: parser.SixteenBit{0},
		OpCode:  parser.StandardQuery,
		Qr:      parser.Response,
		Rcode:   parser.NameErr,
	}
	assert.Equal(t, expectedStruct, h)
}
