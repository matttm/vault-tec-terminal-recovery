package main

import (
	"fmt"

	caravan "caravan"
)

func main() {
	n := 5
	routes := [][]int{
		{0, 1, 80},
		{0, 2, 30},
		{1, 4, 20},
		{2, 3, 40},
		{3, 4, 25},
		{2, 1, 70},
	}

	result := caravan.MinimizeMaxCaravanRisk(n, routes, 0, 4)
	fmt.Printf("Minimized Maximum Route Bottleneck: %d\n", result)
}
