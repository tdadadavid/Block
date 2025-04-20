package main

import "fmt"

func Xor() {
	fileContents := "This is the value of this file"
	checkSum := checksumXor(fileContents)

	// any little change in the file will change the checksum that way we know the files have been tampered with
	downloadedFileContent := "This is the value of this file"
	downloadCheckSum := checksumXor(downloadedFileContent)
	isSame := checkSum == downloadCheckSum
	if isSame {
		fmt.Printf("The files are the same checksum 0x%X\n", checkSum)
	} else {
		fmt.Printf("different checksum. DownloadedFile (0x%X) SentFile (0x%X)\n", checkSum, checksumXor(downloadedFileContent))
	}
}

func checksumXor(fileContents string) byte {
	var checkSum byte = 0

	for i := 0; i < len(fileContents); i++ {
		checkSum ^= fileContents[i]
	}

	return checkSum
}
