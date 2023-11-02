package packets

// Maximum TFTP package size to avoid fragmention - 516 bytes
// TFTP header's size - 4 bytes (Operation code and Block number)
const (
	PackageSize   = 516
	DataBlockSize = PackageSize - 4
)

type OperationCode uint8

const (
	OperationReadRequest OperationCode = iota + 1 // RRQ packet code
	_                                             // WRQ packet code
	OperatonDataResponse                          // DATA packet code
	OperationAcknowledge                          // ACK packet code
	OperationError                                // ERROR packet code
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
