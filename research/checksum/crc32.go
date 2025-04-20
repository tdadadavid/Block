package main

import (
	"fmt"
	"hash/crc32"
)

var POLYNOMIAL uint32 = 0xEDB88320

func CRC32() {
	table := preComputeTable()

	fileContents := "This is the value of this file"

	checksum := Crc32([]byte(fileContents), table)

	goChecksum := crc32.ChecksumIEEE([]byte(fileContents))
	if checksum != goChecksum {
		fmt.Printf("checksums are different 0x%X != 0x%X\n", checksum, goChecksum)
	} else {
		fmt.Printf("checksums are the same 0x%X == 0x%X\n", checksum, goChecksum)
	}
}

func Crc32(val []byte, table [256]uint32) uint32 {
	crc := uint32(0xFFFFFFFF)
	// for each byte in the data
	for _, data := range val {
		// XOR byte with the least significant byte of the CRC
		index := data ^ byte(crc)
		// Shift the CRC to the right and XOR using the pre-computed CRC32 table
		crc = (crc >> 8) ^ table[index]
	}
	// invert the bytes
	return crc ^ 0xFFFFFFFF
}

// This is used to get the CRC value for every possible bit
// This speeds up CRC calculation
func preComputeTable() [256]uint32 {
	var table [256]uint32
	for i := 0; i < 256; i++ {
		crc := uint32(i)
		// this is to simulate division by the polynomial
		for j := 0; j < 8; j++ {
			// if least significant bit set then XOR with polynomial
			if crc&1 == 1 {
				// shift right then XOR
				crc = (crc >> 1) ^ POLYNOMIAL
			} else {
				// shift right
				crc >>= 1
			}
		}
		table[i] = crc
	}
	return table
}
