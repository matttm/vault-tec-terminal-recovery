package main

import (
	"fmt"

	bitsiphon "bitsiphon"
)

func main() {
	packets := []uint32{
		(uint32(0x59) << 24) | 0x00ABCDEF,
		(uint32(0x41) << 24) | 0x00123456,
		(uint32(0x4C) << 24) | 0x00778899,
		(uint32(0x56) << 24) | 0x00000000,
		(uint32(0x55) << 24) | 0x00FF00FF,
		(uint32(0x55) << 24) | 0x00AAAAAA,
	}
	decryptedMessage := make([]byte, len(packets))

	bytesExtracted := bitsiphon.SiphonIntelStream(packets, decryptedMessage)
	finalPayload := decryptedMessage[:bytesExtracted]

	fmt.Printf("Total Bytes Recovered: %d\n", bytesExtracted)
	fmt.Printf("Decrypted Secret Message: %s\n", string(finalPayload))
}
