For this next system design puzzle, we are heading out of the Mojave and straight into the irradiated Capital Wasteland of **Fallout 3**.

You are at the Citadel, working with the Brotherhood of Steel. They’ve salvaged an ancient mainframe from the ruins of **Raven Rock** (the Enclave’s base), but the filesystem is a chaotic, fragmented mess.

---

## The Scenario: The Enclave's Distributed Memory Core

The Enclave used a custom graph-based file storage system. Instead of directories, individual data blocks (Nodes) point directly to other data blocks via logical reference links (Edges).

To prevent the Brotherhood from cloning their AI consciousness profiles, the Enclave implemented a **self-destruct cycle detector** into the hardware registers.

When you boot the core, it passes an initial array of memory linkages. If there is a **directed cycle** anywhere in the active memory sector (meaning a data block loops back to point to an ancestor block, creating an infinite processing loop), the mainframe detects the recursion, triggers a thermal meltdown, and fries the terminal.

Your job is to write a bulletproof execution kernel that detects if a memory sector contains a cycle. If it does, you must abort the boot sequence before the hardware overheats.

---

## The Graph Specification: Memory Sectors

You are given:

* `numBlocks`: The total number of memory blocks, labeled from `0` to `numBlocks - 1`.
* `links`: A 2D array where `links[i] = [u, v]` represents a **directed link** pointing from memory block `u` directly to memory block `v`.

### The Low-Level Twist: Absolute Iterative Safety

The Brotherhood’s terminal runs on an ancient processor with a tiny stack frame. **You cannot use standard recursion (DFS)** to find the cycle. If the Enclave links a chain of 10,000 blocks together linearly, a recursive function will instantly trigger a stack overflow, causing a hardware crash.

### Your Performance Invariants:

1. **Strict Space Complexity:** $O(V)$ auxiliary space. You cannot allocate a giant $V \times V$ adjacency matrix. You must use a flat tracking array.
2. **Strict Time Complexity:** $O(V + E)$ linear time. You must evaluate the sector in a single pass.
3. **No Recursion:** The solution must be completely iterative.

---

## The Strategy Matrix: Kahn’s Algorithm (Topological Sort BFS)

To solve a directed cycle detection problem completely iteratively without a stack frame, we pivot to **Kahn's Algorithm for Topological Sorting** using an in-degree tracking array and a flat BFS queue.

Here is the mechanical intuition:

1. You count how many links point *into* each block (its **In-Degree**).
2. Any block with an In-Degree of `0` is a dead-end entry point—nothing points to it. You throw all `0` in-degree nodes into a flat BFS queue slice.
3. You pop a node from the queue, simulate deleting it from the graph, and decrement the in-degree count of all the nodes it was pointing to.
4. If any of those target nodes hit an in-degree of `0`, you push them into the queue.

If you successfully process **all** blocks in the system, the graph is a Directed Acyclic Graph (DAG)—Boot Sequence Safe! **If you finish and some blocks still have remaining links pointing to them, they are locked in a mutual dependency cycle.** Melt Down Imminent!

---

## The Go Blueprint

Implement the safety kernel in Go matching this signature:

```go
package main

import "fmt"

// IsEnclaveCoreSafe returns true if the memory core is safe to boot (no cycles).
// It returns false if a cycle is detected, indicating a trap.
func IsEnclaveCoreSafe(numBlocks int, links [][]int) bool {
    // Your flat, iterative Kahn's algorithm logic here
}

```

---

## Challenge Dataset

Test your terminal routine against these two memory sectors recovered from Raven Rock:

```go
func main() {
    // SECTOR A: Standard Distributed AI Profile Modules
    numBlocksA := 4
    linksA := [][]int{
        {0, 1}, // Module 0 points to 1
        {1, 2}, // Module 1 points to 2
        {2, 3}, // Module 2 points to 3
    }
    
    // SECTOR B: Decoy Trap Sector
    numBlocksB := 4
    linksB := [][]int{
        {0, 1}, // 0 points to 1
        {1, 2}, // 1 points to 2
        {2, 3}, // 2 points to 3
        {3, 1}, // TRAP! 3 points back to 1, creating a loop (1 -> 2 -> 3 -> 1)
    }

    fmt.Printf("Sector A Safe to Boot: %t\n", IsEnclaveCoreSafe(numBlocksA, linksA)) // Expected: true
    fmt.Printf("Sector B Safe to Boot: %t\n", IsEnclaveCoreSafe(numBlocksB, linksB)) // Expected: false
}

```

### Questions to Trace Before Coding:

1. Why does a cycle prevent the nodes inside it from ever hitting an in-degree of `0`?
2. How does checking the total count of processed nodes at the very end give away whether a cycle exists?

Let's see you save the Brotherhood's terminal!

---

## Local LeetCode-Style Setup

This directory is set up as a small Go coding challenge. The entry point is:

```go
func IsEnclaveCoreSafe(numBlocks int, links [][]int) bool
```

Use these commands from inside `distributed-memory-core`:

```sh
go test -v ./...
go run ./cmd/corecheck
```

Files:

* `solution.go` contains the starter stub for your implementation.
* `solution_test.go` contains visible test cases for the sample sectors, disconnected cycles, independent DAG components, self-loops, and an empty graph. These tests are expected to fail until you implement the stub.
* `cmd/corecheck/main.go` runs the challenge datasets from this README and prints the safety result for each sector.
