package parser

type (
	QR     int
	OPCODE int
	RCODE  int
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
