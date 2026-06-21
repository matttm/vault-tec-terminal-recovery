package main

import (
	"fmt"

	tacticalsat "tacticalsat"
)

func main() {
	k := 3
	stream := []uint64{
		(uint64(10) << 48) | 0x0000FFFFFFFF,
		(uint64(45) << 48) | 0x0000FFFFFFFF,
		(uint64(20) << 48) | 0x0000FFFFFFFF,
		(uint64(35) << 48) | 0x0000FFFFFFFF,
		(uint64(30) << 48) | 0x0000FFFFFFFF,
		(uint64(50) << 48) | 0x0000FFFFFFFF,
		(uint64(15) << 48) | 0x0000FFFFFFFF,
	}
	out := make([]uint16, len(stream)-k+1)

	tacticalsat.AnalyzeTelemetry(stream, k, out)

	fmt.Printf("Parsed Target Trajectories: %v\n", out)
}
