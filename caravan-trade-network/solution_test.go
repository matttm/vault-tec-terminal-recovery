package caravan

import "testing"

func TestMinimizeMaxCaravanRiskChallengeDataset(t *testing.T) {
	n := 5
	routes := [][]int{
		{0, 1, 80},
		{0, 2, 30},
		{1, 4, 20},
		{2, 3, 40},
		{3, 4, 25},
		{2, 1, 70},
	}

	got := MinimizeMaxCaravanRisk(n, routes, 0, 4)
	if got != 40 {
		t.Fatalf("MinimizeMaxCaravanRisk() = %d, want %d", got, 40)
	}
}

func TestMinimizeMaxCaravanRiskUsesBottleneckNotSum(t *testing.T) {
	n := 4
	routes := [][]int{
		{0, 1, 1},
		{1, 2, 1},
		{2, 3, 100},
		{0, 3, 60},
	}

	got := MinimizeMaxCaravanRisk(n, routes, 0, 3)
	if got != 60 {
		t.Fatalf("MinimizeMaxCaravanRisk() = %d, want %d", got, 60)
	}
}

func TestMinimizeMaxCaravanRiskHandlesSourceEqualsDestination(t *testing.T) {
	routes := [][]int{
		{0, 1, 50},
	}

	got := MinimizeMaxCaravanRisk(2, routes, 0, 0)
	if got != 0 {
		t.Fatalf("MinimizeMaxCaravanRisk() = %d, want %d", got, 0)
	}
}

func TestMinimizeMaxCaravanRiskReturnsNegativeOneWhenUnreachable(t *testing.T) {
	n := 4
	routes := [][]int{
		{0, 1, 20},
		{2, 3, 30},
	}

	got := MinimizeMaxCaravanRisk(n, routes, 0, 3)
	if got != -1 {
		t.Fatalf("MinimizeMaxCaravanRisk() = %d, want %d", got, -1)
	}
}

func TestMinimizeMaxCaravanRiskSingleRoute(t *testing.T) {
	routes := [][]int{
		{0, 1, 17},
	}

	got := MinimizeMaxCaravanRisk(2, routes, 0, 1)
	if got != 17 {
		t.Fatalf("MinimizeMaxCaravanRisk() = %d, want %d", got, 17)
	}
}
