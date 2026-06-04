package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_StandardQuery(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      true,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0x00, 0x01},
		QdCount: sixteenBit{0x00, 0x01},
		AnCount: sixteenBit{0},
		NsCount: sixteenBit{0},
		ArCount: sixteenBit{0},
		OpCode:  standardQuery,
		Qr:      query,
		Rcode:   noErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ResponseWithAllFlags(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{0xAB, 0xCD, 0x87, 0x80, 0x00, 0x02, 0x00, 0x01, 0x00, 0x03, 0x00, 0x04}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      true,
		Tc:      true,
		Rd:      true,
		Ra:      true,
		Z:       0,
		Id:      sixteenBit{0xAB, 0xCD},
		QdCount: sixteenBit{0x00, 0x02},
		AnCount: sixteenBit{0x00, 0x01},
		NsCount: sixteenBit{0x00, 0x03},
		ArCount: sixteenBit{0x00, 0x04},
		OpCode:  standardQuery,
		Qr:      response,
		Rcode:   noErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_InverseQuery(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{0x12, 0x34, 0x08, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0x12, 0x34},
		QdCount: sixteenBit{0x00, 0x01},
		AnCount: sixteenBit{0},
		NsCount: sixteenBit{0},
		ArCount: sixteenBit{0},
		OpCode:  inverseQuery,
		Qr:      query,
		Rcode:   noErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ServerStatus(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{0xFF, 0xFF, 0x90, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0xFF, 0xFF},
		QdCount: sixteenBit{0},
		AnCount: sixteenBit{0},
		NsCount: sixteenBit{0},
		ArCount: sixteenBit{0},
		OpCode:  serverStatus,
		Qr:      response,
		Rcode:   serverFailure,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_VariousRCodes(t *testing.T) {
	tests := []struct {
		name     string
		rcode    uint8
		expected RCODE
	}{
		{"noError", 0, noErrorCon},
		{"formatError", 1, formatErr},
		{"serverFailure", 2, serverFailure},
		{"nameError", 3, nameErr},
		{"notImplemented", 4, notImplemented},
		{"refused", 5, refused},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := createHeaders()
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
	h := createHeaders()
	dnsMock := [12]byte{0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0, 0},
		QdCount: sixteenBit{0xFF, 0xFF},
		AnCount: sixteenBit{0xFF, 0xFF},
		NsCount: sixteenBit{0xFF, 0xFF},
		ArCount: sixteenBit{0xFF, 0xFF},
		OpCode:  standardQuery,
		Qr:      query,
		Rcode:   noErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ZeroBuffer(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0, 0},
		QdCount: sixteenBit{0, 0},
		AnCount: sixteenBit{0, 0},
		NsCount: sixteenBit{0, 0},
		ArCount: sixteenBit{0, 0},
		OpCode:  standardQuery,
		Qr:      query,
		Rcode:   noErrorCon,
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
			h := createHeaders()
			dnsMock := [12]byte{0x00, 0x00, 0x00, tt.val, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			err := h.Parse(dnsMock)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Z must be 0")
		})
	}
}

func TestParser_RealWorldGoogleDNS(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{0x12, 0x34, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      true,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0x12, 0x34},
		QdCount: sixteenBit{0x00, 0x01},
		AnCount: sixteenBit{0},
		NsCount: sixteenBit{0},
		ArCount: sixteenBit{0},
		OpCode:  standardQuery,
		Qr:      query,
		Rcode:   noErrorCon,
	}
	assert.Equal(t, expectedStruct, h)
}

func TestParser_ResponseWithNXDOMAIN(t *testing.T) {
	h := createHeaders()
	dnsMock := [12]byte{0x00, 0x01, 0x80, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}
	err := h.Parse(dnsMock)
	assert.NoError(t, err)
	expectedStruct := &Header{
		Aa:      false,
		Tc:      false,
		Rd:      false,
		Ra:      false,
		Z:       0,
		Id:      sixteenBit{0x00, 0x01},
		QdCount: sixteenBit{0x00, 0x01},
		AnCount: sixteenBit{0},
		NsCount: sixteenBit{0x00, 0x01},
		ArCount: sixteenBit{0},
		OpCode:  standardQuery,
		Qr:      response,
		Rcode:   nameErr,
	}
	assert.Equal(t, expectedStruct, h)
}
