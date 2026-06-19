package vaulttec

import (
	"errors"
	"unsafe"
)

const (
	headerLen         = 7
	metadataOffset    = 4
	payloadLenOffset  = 6
	payloadOffset     = 7
	requiredClearance = 0x0F
	requiredSector    = 0x65
	xorEncryption     = 0x01
	xorKey            = 0x55
)

var (
	ErrPacketTooShort        = errors.New("packet too short")
	ErrInvalidPreamble       = errors.New("invalid preamble")
	ErrUnauthorized          = errors.New("unauthorized packet")
	ErrUnsupportedEncryption = errors.New("unsupported encryption type")
	ErrTruncatedPayload      = errors.New("truncated payload")
)

// DecodeVaultOverride parses the raw stream in-place.
// It returns the decoded passcode string, or an error if the packet is corrupt or unauthorized.
// On success, the returned string aliases stream's payload bytes to satisfy the zero-allocation constraint.
func DecodeVaultOverride(stream []byte) (string, error) {
	if len(stream) < headerLen {
		return "", ErrPacketTooShort
	}

	if stream[0] != 0x52 || stream[1] != 0x4F || stream[2] != 0x42 || stream[3] != 0x43 {
		return "", ErrInvalidPreamble
	}

	metadata := uint16(stream[metadataOffset])<<8 | uint16(stream[metadataOffset+1])
	clearance := metadata & 0x000F
	sector := (metadata >> 4) & 0x00FF
	encryption := (metadata >> 12) & 0x000F

	if clearance != requiredClearance || sector != requiredSector {
		return "", ErrUnauthorized
	}

	payloadLen := int(stream[payloadLenOffset])
	if len(stream)-payloadOffset < payloadLen {
		return "", ErrTruncatedPayload
	}

	payload := stream[payloadOffset : payloadOffset+payloadLen]
	switch encryption {
	case 0:
	case xorEncryption:
		for i := range payload {
			payload[i] ^= xorKey
		}
	default:
		return "", ErrUnsupportedEncryption
	}

	return stringView(payload), nil
}

func stringView(payload []byte) string {
	if len(payload) == 0 {
		return ""
	}

	return unsafe.String(&payload[0], len(payload))
}
