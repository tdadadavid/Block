package toolkit

import (
	"bytes"
	"encoding/binary"
)

// Serialize a string (prefix length, then data)
func SerializeString(buf *bytes.Buffer, str string) error {
	strLen := uint32(len(str))
	if err := binary.Write(buf, binary.LittleEndian, strLen); err != nil {
		return err
	}
	_, err := buf.Write([]byte(str))
	return err
}

// Deserialize a string (read length first, then data)
func DeserializeString(buf *bytes.Reader) (string, error) {
	var strLen uint32
	if err := binary.Read(buf, binary.LittleEndian, &strLen); err != nil {
		return "", err
	}

	strBytes := make([]byte, strLen)
	if _, err := buf.Read(strBytes); err != nil {
		return "", err
	}
	return string(strBytes), nil
}