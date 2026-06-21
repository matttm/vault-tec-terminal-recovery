package vaulttec

import (
	"encoding/binary"
	"errors"
)

// DecodeVaultOverride parses the raw stream in-place.
// It returns the decoded passcode string, or an error if the packet is corrupt or unauthorized.
func DecodeVaultOverride(stream []byte) (string, error) {
	if len(stream) < 7 {
		return "", errors.New("packet too short")
	}
	if !(stream[0] == 0x52 && stream[1] == 0x4F && stream[2] == 0x42 && stream[3] == 0x43) {
		return "", errors.New("invalid preamble")
	}
	meta := binary.BigEndian.Uint16(stream[4:6]) // 0x165F
	clearance := meta & 0x000F        // bits 0-3  => 0xF
	secId := (meta >> 4) & 0x00FF     // bits 4-11 => 0x65
	encType := (meta >> 12) & 0x000F  // bits 12-15 => 0x1
	// log.Printf("clearance: %d, secId: %d, encType: %d", clearance, secId, encType)
	if clearance != 0xF {
		return "", errors.New("insufficient clearance")
	}
	if secId != 0x65 {
		return "", errors.New("unauthorized security ID")
	}
	plen := stream[6] // 0x08
	if len(stream) < int(7+plen) {
		return "", errors.New("packet too short for payload")
	}
	payload := stream[7 : 7+plen] // "12345678"
	if encType == 0x1 {
		// decrypt in place
		for i := range payload {
			payload[i] ^= 0x55
		}
	} else if encType != 0x0 {
		return "", errors.New("unsupported encryption type")
	}
	return string(payload), nil
}
