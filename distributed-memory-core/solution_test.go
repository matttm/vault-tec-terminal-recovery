package memorycore

import "testing"

func TestIsEnclaveCoreSafeChallengeSafeSector(t *testing.T) {
	numBlocks := 4
	// Graph:
	// 0 -> 1 -> 2 -> 3
	links := [][]int{
		{0, 1},
		{1, 2},
		{2, 3},
	}

	got := IsEnclaveCoreSafe(numBlocks, links)
	if got != true {
		t.Fatalf("IsEnclaveCoreSafe() = %t, want true", got)
	}
}

func TestIsEnclaveCoreSafeChallengeTrapSector(t *testing.T) {
	numBlocks := 4
	// Graph:
	// 0 -> 1 -> 2 -> 3
	//      ^         |
	//      |_________|
	links := [][]int{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 1},
	}

	got := IsEnclaveCoreSafe(numBlocks, links)
	if got != false {
		t.Fatalf("IsEnclaveCoreSafe() = %t, want false", got)
	}
}

func TestIsEnclaveCoreSafeDetectsDisconnectedCycle(t *testing.T) {
	numBlocks := 6
	// Graph:
	// 0 -> 1 -> 2
	//
	// 3 -> 4
	// ^    |
	// |____|
	//
	// 5 is isolated.
	links := [][]int{
		{0, 1},
		{1, 2},
		{3, 4},
		{4, 3},
	}

	got := IsEnclaveCoreSafe(numBlocks, links)
	if got != false {
		t.Fatalf("IsEnclaveCoreSafe() = %t, want false", got)
	}
}

func TestIsEnclaveCoreSafeHandlesIndependentDAGComponents(t *testing.T) {
	numBlocks := 6
	// Graph:
	// 0 -> 2 <- 1
	//
	//      /-> 4
	// 3 --|
	//      \-> 5
	links := [][]int{
		{0, 2},
		{1, 2},
		{3, 4},
		{3, 5},
	}

	got := IsEnclaveCoreSafe(numBlocks, links)
	if got != true {
		t.Fatalf("IsEnclaveCoreSafe() = %t, want true", got)
	}
}

func TestIsEnclaveCoreSafeDetectsSelfLoop(t *testing.T) {
	numBlocks := 3
	// Graph:
	// 0 -> 1
	//
	// 2 loops back to itself.
	links := [][]int{
		{0, 1},
		{2, 2},
	}

	got := IsEnclaveCoreSafe(numBlocks, links)
	if got != false {
		t.Fatalf("IsEnclaveCoreSafe() = %t, want false", got)
	}
}

func TestIsEnclaveCoreSafeHandlesEmptyGraph(t *testing.T) {
	// Graph:
	// No blocks and no links.
	got := IsEnclaveCoreSafe(0, nil)
	if got != true {
		t.Fatalf("IsEnclaveCoreSafe() = %t, want true", got)
	}
}
