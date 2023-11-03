package packets

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
)

// DATA packet: Opeartion Code | Block Number | Payload
// DATA packet's size: 2 bytes + 2 bytes + n bytes

type DataPacket struct {
	BlockNumber uint16
	Payload     io.Reader

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func (dp *DataPacket) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(PacketSize)

	dp.BlockNumber++

	err := binary.Write(buf, binary.BigEndian, OperationDataResponse)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, dp.BlockNumber)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(buf, dp.Payload, DataBlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (dp *DataPacket) UnmarshalBinary(data []byte) error {
	if l := len(data); l < 4 || l > PacketSize {
		return errors.New("invalid DATA packet, length of a packet is incorrect")
	}

	buf := bytes.NewBuffer(data)

	var code OperationCode
	err := binary.Read(buf, binary.BigEndian, &code)
	if err != nil || code != OperationDataResponse {
		return errors.New("invalid DATA packet, packet operation code is incorrect")
	}

	err = binary.Read(buf, binary.BigEndian, &dp.BlockNumber)
	if err != nil {
		return errors.New("invalid DATA packet, unable to read packet's number")
	}

	err = binary.Read(buf, binary.BigEndian, dp.Payload)
	if err != nil {
		return errors.New("invalid DATA packet, unable to read the payload")
	}

	// dp.Payload = bytes.NewBuffer(data[4:])

	return nil
}
