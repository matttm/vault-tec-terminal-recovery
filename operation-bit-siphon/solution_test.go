package bitsiphon

import (
	"bytes"
	"reflect"
	"testing"
)

func encryptedPacket(hidden byte, lowBits uint32) uint32 {
	return uint32(hidden)<<24 | (lowBits & 0x00FFFFFF)
}

func TestSiphonIntelStreamChallengeDataset(t *testing.T) {
	packets := []uint32{
		encryptedPacket(0x59, 0x00ABCDEF),
		encryptedPacket(0x41, 0x00123456),
		encryptedPacket(0x4C, 0x00778899),
		encryptedPacket(0x56, 0x00000000),
		encryptedPacket(0x55, 0x00FF00FF),
		encryptedPacket(0x55, 0x00AAAAAA),
	}
	decryptedMessage := make([]byte, len(packets))

	got := SiphonIntelStream(packets, decryptedMessage)

	if got != 6 {
		t.Fatalf("SiphonIntelStream() = %d, want 6", got)
	}
	if string(decryptedMessage[:got]) != "XCORPS" {
		t.Fatalf("SiphonIntelStream() wrote %q, want %q", decryptedMessage[:got], "XCORPS")
	}
}

func TestSiphonIntelStreamExtractsOnlyTopPaddingByte(t *testing.T) {
	packets := []uint32{
		encryptedPacket('A'^0x01, 0x00000000),
		encryptedPacket('B'^0x02, 0x00FFFFFF),
		encryptedPacket('C'^0x03, 0x0055AA55),
	}
	decryptedMessage := make([]byte, len(packets))

	got := SiphonIntelStream(packets, decryptedMessage)

	if got != 3 {
		t.Fatalf("SiphonIntelStream() = %d, want 3", got)
	}
	if string(decryptedMessage[:got]) != "ABC" {
		t.Fatalf("SiphonIntelStream() wrote %q, want %q", decryptedMessage[:got], "ABC")
	}
}

func TestSiphonIntelStreamStopsAtOutputCapacity(t *testing.T) {
	packets := []uint32{
		encryptedPacket('O'^0x01, 0x00000000),
		encryptedPacket('K'^0x02, 0x00000000),
		encryptedPacket('X'^0x03, 0x00000000),
	}
	decryptedMessage := []byte{0, 0}

	got := SiphonIntelStream(packets, decryptedMessage)

	if got != 2 {
		t.Fatalf("SiphonIntelStream() = %d, want 2", got)
	}
	if string(decryptedMessage) != "OK" {
		t.Fatalf("SiphonIntelStream() wrote %q, want %q", decryptedMessage, "OK")
	}
}

func TestSiphonIntelStreamLeavesOutputAloneForEmptyInput(t *testing.T) {
	decryptedMessage := []byte{'x', 'y', 'z'}
	before := append([]byte(nil), decryptedMessage...)

	got := SiphonIntelStream(nil, decryptedMessage)

	if got != 0 {
		t.Fatalf("SiphonIntelStream() = %d, want 0", got)
	}
	if !reflect.DeepEqual(decryptedMessage, before) {
		t.Fatalf("SiphonIntelStream() changed output to %v, want unchanged %v", decryptedMessage, before)
	}
}

func TestSiphonIntelStreamDoesNotAllocateOnSuccessfulSiphon(t *testing.T) {
	packets := []uint32{
		encryptedPacket('V'^0x01, 0x00ABCDEF),
		encryptedPacket('A'^0x02, 0x00123456),
		encryptedPacket('U'^0x03, 0x00778899),
		encryptedPacket('L'^0x04, 0x00000000),
		encryptedPacket('T'^0x05, 0x00FF00FF),
	}
	decryptedMessage := make([]byte, len(packets))
	want := []byte("VAULT")

	got := SiphonIntelStream(packets, decryptedMessage)
	if got != 5 || !bytes.Equal(decryptedMessage[:got], want) {
		t.Fatalf("SiphonIntelStream() wrote %q with count %d, want %q with count 5", decryptedMessage[:got], got, "VAULT")
	}

	allocs := testing.AllocsPerRun(1000, func() {
		clear(decryptedMessage)
		got = SiphonIntelStream(packets, decryptedMessage)
		if got != 5 || !bytes.Equal(decryptedMessage[:got], want) {
			panic("wrong decrypted output")
		}
	})

	if allocs != 0 {
		t.Fatalf("SiphonIntelStream allocated %.0f times, want zero", allocs)
	}
}
