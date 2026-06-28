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

func TestIsEnclaveCoreSafeDetectsFiveCyclesAcrossSixtyNodes(t *testing.T) {
	numBlocks := 60
	// Graph:
	// Five cyclic components:
	// 0 -> 1 -> 2 -> 0
	// 10 -> 11 -> 12 -> 13 -> 10
	// 20 -> 21 -> 20
	// 30 -> 31 -> 32 -> 33 -> 34 -> 30
	// 45 -> 46 -> 47 -> 48 -> 49 -> 45
	//
	// All other nodes are part of acyclic chains between those components.
	links := [][]int{
		{0, 1},
		{1, 2},
		{2, 0},
		{3, 4},
		{4, 5},
		{5, 6},
		{6, 7},
		{7, 8},
		{8, 9},
		{10, 11},
		{11, 12},
		{12, 13},
		{13, 10},
		{14, 15},
		{15, 16},
		{16, 17},
		{17, 18},
		{18, 19},
		{20, 21},
		{21, 20},
		{22, 23},
		{23, 24},
		{24, 25},
		{25, 26},
		{26, 27},
		{27, 28},
		{28, 29},
		{30, 31},
		{31, 32},
		{32, 33},
		{33, 34},
		{34, 30},
		{35, 36},
		{36, 37},
		{37, 38},
		{38, 39},
		{39, 40},
		{40, 41},
		{41, 42},
		{42, 43},
		{43, 44},
		{45, 46},
		{46, 47},
		{47, 48},
		{48, 49},
		{49, 45},
		{50, 51},
		{51, 52},
		{52, 53},
		{53, 54},
		{54, 55},
		{55, 56},
		{56, 57},
		{57, 58},
		{58, 59},
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
