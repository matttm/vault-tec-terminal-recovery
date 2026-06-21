package tacticalsat

import (
	"reflect"
	"testing"
)

func packTelemetry(rad uint16) uint64 {
	return uint64(rad)<<48 | 0x0000FFFFFFFF
}

func TestAnalyzeTelemetryChallengeDataset(t *testing.T) {
	stream := []uint64{
		packTelemetry(10),
		packTelemetry(45),
		packTelemetry(20),
		packTelemetry(35),
		packTelemetry(30),
		packTelemetry(50),
		packTelemetry(15),
	}
	out := make([]uint16, len(stream)-3+1)

	AnalyzeTelemetry(stream, 3, out)

	want := []uint16{45, 45, 35, 50, 50}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("AnalyzeTelemetry() wrote %v, want %v", out, want)
	}
}

func TestAnalyzeTelemetryExtractsOnlyRadCountBits(t *testing.T) {
	stream := []uint64{
		uint64(1)<<48 | 0x000000000000,
		uint64(255)<<48 | 0x123456789ABC,
		uint64(42)<<48 | 0xFFFFFFFFFFFF,
	}
	out := make([]uint16, len(stream))

	AnalyzeTelemetry(stream, 1, out)

	want := []uint16{1, 255, 42}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("AnalyzeTelemetry() wrote %v, want %v", out, want)
	}
}

func TestAnalyzeTelemetryHandlesWholeStreamWindow(t *testing.T) {
	stream := []uint64{
		packTelemetry(7),
		packTelemetry(4),
		packTelemetry(99),
		packTelemetry(12),
	}
	out := make([]uint16, 1)

	AnalyzeTelemetry(stream, len(stream), out)

	want := []uint16{99}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("AnalyzeTelemetry() wrote %v, want %v", out, want)
	}
}

func TestAnalyzeTelemetryLeavesOutputAloneForInvalidInputs(t *testing.T) {
	tests := []struct {
		name      string
		telemetry []uint64
		k         int
		out       []uint16
	}{
		{
			name:      "empty telemetry",
			telemetry: nil,
			k:         3,
			out:       []uint16{9, 9},
		},
		{
			name:      "zero window",
			telemetry: []uint64{packTelemetry(1), packTelemetry(2)},
			k:         0,
			out:       []uint16{9, 9},
		},
		{
			name:      "window larger than telemetry",
			telemetry: []uint64{packTelemetry(1), packTelemetry(2)},
			k:         3,
			out:       []uint16{9, 9},
		},
		{
			name:      "output too short",
			telemetry: []uint64{packTelemetry(1), packTelemetry(2), packTelemetry(3)},
			k:         2,
			out:       []uint16{9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := append([]uint16(nil), tt.out...)
			AnalyzeTelemetry(tt.telemetry, tt.k, tt.out)
			if !reflect.DeepEqual(tt.out, before) {
				t.Fatalf("AnalyzeTelemetry() changed out to %v, want unchanged %v", tt.out, before)
			}
		})
	}
}

func TestAnalyzeTelemetryDoesNotAllocateOnSuccessfulAnalysis(t *testing.T) {
	stream := []uint64{
		packTelemetry(10),
		packTelemetry(45),
		packTelemetry(20),
		packTelemetry(35),
		packTelemetry(30),
		packTelemetry(50),
		packTelemetry(15),
	}
	out := make([]uint16, len(stream)-3+1)
	want := []uint16{45, 45, 35, 50, 50}

	AnalyzeTelemetry(stream, 3, out)
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("AnalyzeTelemetry() wrote %v before allocation check, want %v", out, want)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		clear(out)
		AnalyzeTelemetry(stream, 3, out)
	})

	if allocs != 0 {
		t.Fatalf("AnalyzeTelemetry allocated %.0f times, want zero", allocs)
	}
}
