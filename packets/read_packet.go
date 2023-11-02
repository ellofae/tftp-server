package packets

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// READ package: Operation code | Filename | 0 byte | Mode | 0 byte |
// READ package size: 2 bytes + n bytes + 1 byte + n bytes + 1 byte

type ReadRequest struct {
	Filename string // Name of a requested file
	Mode     string // Line-ending format (octet, netascii)

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func (rq ReadRequest) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if rq.Mode != "" {
		mode = strings.ToLower(rq.Mode)
	}

	bufCapacity := len(rq.Filename) + len(rq.Mode) + 4

	buf := new(bytes.Buffer)
	buf.Grow(bufCapacity)

	err := binary.Write(buf, binary.BigEndian, OperationReadRequest)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(rq.Filename)
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(mode)
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (rq *ReadRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	var code OperationCode
	err := binary.Read(buf, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OperationReadRequest {
		return errors.New("invalid operation code for the READ packet")
	}

	// Reading bytes from buffer until 0 delimiter, including the delimiter
	rq.Filename, err = buf.ReadString(0)
	if err != nil {
		return fmt.Errorf("invalid READ packet, error: %w", err)
	}

	// Removing the 0-byte at the end of a filename string
	rq.Filename = strings.TrimRight(rq.Filename, "\x00")
	if len(rq.Filename) == 0 {
		return errors.New("invalid READ packet due to incorrect filename")
	}

	rq.Mode, err = buf.ReadString(0)
	if err != nil {
		return fmt.Errorf("invalid READ packet, error: %w", err)
	}

	rq.Mode = strings.TrimRight(rq.Mode, "\x00")
	if len(rq.Mode) == 0 {
		return errors.New("invalid READ packet due to incorrect mode format")
	}

	actual := strings.ToLower(rq.Mode)
	if actual != "octet" {
		return errors.New("invalid READ packet, only binary transfers are supported")
	}

	return nil
}
