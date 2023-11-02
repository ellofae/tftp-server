package packets

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
)

// Acknowledgment packets are 4 bytes long
// ACK packet: Operation Code | Block Number

type AckPacket struct {
	BlockNumber uint16 // uint16: 2 bytes for amount of packages => 516 bytes * 65535 ~= 33.8 MB file's max size

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func (ap AckPacket) MarshalBinary() ([]byte, error) {
	bufCapacity := 4

	buf := new(bytes.Buffer)
	buf.Grow(bufCapacity)

	err := binary.Write(buf, binary.BigEndian, OperationAcknowledge)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, ap.BlockNumber)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ap *AckPacket) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	var code OperationCode
	err := binary.Read(buf, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OperationAcknowledge {
		return errors.New("invalid ACKS packet, operation code is not Acknowledgment code")
	}

	return binary.Read(buf, binary.BigEndian, ap.BlockNumber)
}
