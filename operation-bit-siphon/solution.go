package bitsiphon

import "log"

// SiphonIntelStream extracts, decrypts, and writes hidden data from the packet
// stream into the pre-allocated decryptedMessage buffer. It returns the total
// number of secret bytes safely extracted.
func SiphonIntelStream(packets []uint32, decryptedMessage []byte) int {
	outputCapacity := len(decryptedMessage)
	bitMask := 0x01
	for i := range packets {
		if i >= outputCapacity {
			log.Printf("Output buffer full: processed %d packets, output capacity %d", i, outputCapacity)
			return outputCapacity
		}
		padding := (packets[i] >> 24) & 0xFF
		decryptedMessage[i] = byte(padding) ^ byte(bitMask)
		bitMask++
		log.Printf("Packet %d: padding=0x%02X, decrypted=0x%02X", i, padding, decryptedMessage[i])
	}
	return len(packets)
}
