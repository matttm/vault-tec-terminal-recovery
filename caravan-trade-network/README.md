Let’s shift gears from raw memory tricks to a classic **Graph Theory Optimization** problem wrapped in the *Fallout* wasteland lore.

This one tests your ability to model network connectivity, track state cleanly, and hit a tight time complexity without letting your memory allocation spin out of control.

---

## The Scenario: The Caravan Trade Network

You are a software engineer hired by the **Crimson Caravan Company** in the Mojave Wasteland. You need to optimize the trade routes connecting various settlements (Goodsprings, Novac, Primm, etc.).

The wasteland is dangerous, so caravans can only travel along specific, scouted routes. You are given an array of these routes, where each route connects two settlements and has a risk rating representing the worst expected threat on that segment: raider ambushes, radiation pockets, minefields, or hostile patrols.

The Caravan Master wants to set up a supply system between a starting settlement (`source`) and a destination settlement (`destination`). The caravan can restock, repair, and recover at every settlement it reaches, so total hardship over the whole journey is not the deciding factor. The real failure point is the most dangerous single route segment: if one leg exceeds the escort crew's threat rating, the caravan is lost before the next checkpoint.

To choose a survivable route, you must find a path where the **maximum risk of any single route on that path is minimized**.

---

## The Core Constraint & Twist

You are given:

* `n`: The total number of settlements (numbered `0` to `n-1`).
* `routes`: A 2D array where `routes[i] = [u, v, risk]` means there is a bidirectional path between settlement `u` and settlement `v` with a specific `risk` factor.
* `source`: The starting settlement index.
* `destination`: The ending settlement index.

### The Twist: The Speed Trap

The Crimson Caravan needs this calculated instantly for hundreds of dynamic raider movements.

* A brute-force Depth-First Search (DFS) exploring all paths will cause a stack overflow and fail the time constraints.
* A standard shortest-path Dijkstra tracks *cumulative* distance by adding edges together. Here, settlements act as safe checkpoints, so we do not care about the *sum* of the risks; we only care about the **maximum bottleneck edge** on the path.

### Your Performance Invariants:

1. **Time Complexity:** Must run better than $O(V^2)$ where $V$ is the number of settlements. Aim for $O(E \log V)$ or $O(E \log W)$ where $E$ is the number of routes and $W$ is the maximum risk range.
2. **Space Complexity:** $O(V + E)$ to store the graph network representation. No memory leaks or massive redundant tracking structures.

---

## The Go Blueprint

Implement the optimization engine in Go matching this signature:

```go
package main

import "fmt"

// MinimizeMaxCaravanRisk finds a path from source to destination 
// that minimizes the maximum risk score of any single route along that path.
func MinimizeMaxCaravanRisk(n int, routes [][]int, source int, destination int) int {
    // Your graph optimization logic here
}

```

---

## Challenge Dataset

Test your tracking logic against this wasteland route map:

```go
func main() {
    n := 5
    // routes[i] = [settlement_A, settlement_B, risk_level]
    routes := [][]int{
        {0, 1, 80}, // Goodsprings to Primm (High raider risk)
        {0, 2, 30}, // Goodsprings to Outpost A (Low risk)
        {1, 4, 20}, // Primm to Vegas (Very low risk)
        {2, 3, 40}, // Outpost A to Novac (Medium risk)
        {3, 4, 25}, // Novac to Vegas (Low risk)
        {2, 1, 70}, // Outpost A to Primm (High risk)
    }
    
    source := 0       // Goodsprings
    destination := 4  // Vegas

    result := MinimizeMaxCaravanRisk(n, routes, source, destination)
    
    // Expected Output: 40
    // Path: 0 -> 2 -> 3 -> 4 
    // Edges crossed have risks: 30, 40, 25. The maximum bottleneck is 40.
    // (Taking 0 -> 1 -> 4 would have an edge of 80, which is worse!)
    fmt.Printf("Minimized Maximum Route Bottleneck: %d\n", result)
}

```

---

### Strategy Hints Before You Code:

There are two primary ways advanced systems engineers crack this puzzle:

1. **The Modified Dijkstra (Max-Min Priority Queue):** Instead of updating a distance array with `dist[u] + weight`, you maintain a `max_risk_seen` array, updating it with $\max(\text{current\_risk}, \text{edge\_weight})$ using a Min-Heap.
2. **Binary Search + BFS/DFS:** The risk levels have a fixed range (e.g., $0$ to $100$ or up to the max edge weight). Can you binary search for the perfect "maximum allowed risk threshold," and use a simple, fast linear BFS to check if a valid path exists using *only* routes below that threshold?

Which pattern sounds like the cleanest architecture to map out first?

---

## Local LeetCode-Style Setup

This directory is set up as a small Go coding challenge. The entry point is:

```go
func MinimizeMaxCaravanRisk(n int, routes [][]int, source int, destination int) int
```

Use these commands from inside `caravan-trade-network`:

```sh
go test -v ./...
go run ./cmd/caravan
```

Files:

* `solution.go` contains the starter stub for your implementation.
* `solution_test.go` contains visible test cases for the sample dataset, bottleneck-vs-sum behavior, same-source destination, unreachable destination, and a single-route graph. These tests are expected to fail until you implement the stub.
* `cmd/caravan/main.go` runs the challenge dataset from this README and prints the result.

For this local setup, return `-1` when `destination` is unreachable from `source`.
