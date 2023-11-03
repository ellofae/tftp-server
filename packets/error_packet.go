package packets

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"strings"
)

// ERROR packet: Operation Code | Error code | Message | 0-byte
// ERROR packet's size: 2 bytes + 2 bytes + n bytes + 1 byte

type ErrorPacket struct {
	Error   ErrorCode
	Message string

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func (rp ErrorPacket) MarshalBinary() ([]byte, error) {
	bufCapacity := len(rp.Message) + 5

	buf := new(bytes.Buffer)
	buf.Grow(bufCapacity)

	err := binary.Write(buf, binary.BigEndian, OperationError)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, rp.Error)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(rp.Message)
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (rp *ErrorPacket) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	var code OperationCode
	err := binary.Read(buf, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OperationError {
		return errors.New("invalid ERROR packet, operation code is not Error(OperationError)")
	}

	err = binary.Read(buf, binary.BigEndian, &rp.Error)
	if err != nil {
		return err
	}

	rp.Message, err = buf.ReadString(0)
	if err != nil {
		return err
	}

	rp.Message = strings.TrimRight(rp.Message, "\x00")

	return nil
}
