package main

import (
	"fmt"

	memorycore "memorycore"
)

func main() {
	numBlocksA := 4
	linksA := [][]int{
		{0, 1},
		{1, 2},
		{2, 3},
	}

	numBlocksB := 4
	linksB := [][]int{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 1},
	}

	fmt.Printf("Sector A Safe to Boot: %t\n", memorycore.IsEnclaveCoreSafe(numBlocksA, linksA))
	fmt.Printf("Sector B Safe to Boot: %t\n", memorycore.IsEnclaveCoreSafe(numBlocksB, linksB))
}
