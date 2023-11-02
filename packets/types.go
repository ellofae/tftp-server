package packets

// Maximum TFTP package size to avoid fragmention - 516 bytes
// TFTP header's size - 4 bytes (Operation code and Block number)
const (
	PacketSize    = 516
	DataBlockSize = PacketSize - 4
)

type OperationCode uint8

const (
	OperationReadRequest  OperationCode = iota + 1 // RRQ packet code
	_                                              // WRQ packet code
	OperationDataResponse                          // DATA packet code
	OperationAcknowledge                           // ACK packet code
	OperationError                                 // ERROR packet code
)

type ErrorCode uint8

// TFTP (RFC 1350) errors
const (
	ErrUnknown ErrorCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)
