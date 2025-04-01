package toolkit

import (
	"bytes"
	"encoding/binary"
)

// SerializeString a string (prefix length, then data)
func SerializeString(buf *bytes.Buffer, str string) (err error) {
	strLen := uint32(len(str))
	if err := binary.Write(buf, binary.LittleEndian, strLen); err != nil {
		return err
	}
	_, err = buf.Write([]byte(str))
	return err
}

// DeserializeString a string (read length first, then data)
func DeserializeString(buf *bytes.Reader) (val string, err error) {
	var strLen uint32
	if err = binary.Read(buf, binary.LittleEndian, &strLen); err != nil {
		return val, err
	}

	strBytes := make([]byte, strLen)
	if _, err := buf.Read(strBytes); err != nil {
		return val, err
	}
	val = string(strBytes)
	return val, err
}
