package vaulttec

import (
	"testing"
)

var recoveredStream = [...]byte{
	0x52, 0x4F, 0x42, 0x43,
	0x16, 0x5F,
	0x0C,
	0x13, 0x34, 0x37, 0x37,
	0x3A, 0x20, 0x13, 0x3A,
	0x3F, 0x39, 0x33, 0x23,
}

func TestDecodeVaultOverrideRecoveredStream(t *testing.T) {
	stream := append([]byte(nil), recoveredStream[:]...)

	got, err := DecodeVaultOverride(stream)
	if err != nil {
		t.Fatalf("DecodeVaultOverride returned error: %v", err)
	}
	if got != "FabbouFojlfv" {
		t.Fatalf("DecodeVaultOverride() = %q, want %q", got, "FabbouFojlfv")
	}
	if string(stream[7:]) != got {
		t.Fatalf("payload was not decrypted in place: %q", stream[7:])
	}
}

func TestDecodeVaultOverridePlaintextPacket(t *testing.T) {
	stream := []byte{
		0x52, 0x4F, 0x42, 0x43,
		0x06, 0x5F,
		0x05,
		'V', 'A', 'U', 'L', 'T',
	}

	got, err := DecodeVaultOverride(stream)
	if err != nil {
		t.Fatalf("DecodeVaultOverride returned error: %v", err)
	}
	if got != "VAULT" {
		t.Fatalf("DecodeVaultOverride() = %q, want %q", got, "VAULT")
	}
}

func TestDecodeVaultOverrideRejectsCorruptPackets(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{
			name: "short header",
			in:   []byte{0x52, 0x4F, 0x42},
		},
		{
			name: "invalid preamble",
			in:   []byte{0x52, 0x4F, 0x42, 0x58, 0x16, 0x5F, 0x00},
		},
		{
			name: "wrong clearance",
			in:   []byte{0x52, 0x4F, 0x42, 0x43, 0x16, 0x5E, 0x00},
		},
		{
			name: "wrong sector",
			in:   []byte{0x52, 0x4F, 0x42, 0x43, 0x16, 0x4F, 0x00},
		},
		{
			name: "unsupported encryption",
			in:   []byte{0x52, 0x4F, 0x42, 0x43, 0x26, 0x5F, 0x00},
		},
		{
			name: "truncated payload",
			in:   []byte{0x52, 0x4F, 0x42, 0x43, 0x16, 0x5F, 0x03, 0x13},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeVaultOverride(tt.in)
			if err == nil {
				t.Fatal("DecodeVaultOverride() error = nil, want non-nil")
			}
		})
	}
}

func TestDecodeVaultOverrideOnlyAllocatesFinalStringOnSuccessfulDecode(t *testing.T) {
	stream := make([]byte, len(recoveredStream))
	copy(stream, recoveredStream[:])

	got, err := DecodeVaultOverride(stream)
	if err != nil {
		t.Fatalf("DecodeVaultOverride returned error before allocation check: %v", err)
	}
	if got != "FabbouFojlfv" {
		t.Fatalf("DecodeVaultOverride() = %q, want %q", got, "FabbouFojlfv")
	}

	allocs := testing.AllocsPerRun(1000, func() {
		copy(stream, recoveredStream[:])
		got, err = DecodeVaultOverride(stream)
		if err != nil {
			panic(err)
		}
		if got != "FabbouFojlfv" {
			panic(got)
		}
	})

	// The function should only allocate once for the final string, and not for any intermediate buffers.
	if allocs > 1 {
		t.Fatalf("DecodeVaultOverride allocated %.0f times, want at most one final string allocation", allocs)
	}
}
