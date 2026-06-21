Here is a low-level systems puzzle that pairs **strict time complexity constraints** with an interesting architecture twist: a **Lock-Free Ring Buffer Log** salvaged from a military defense satellite network in the *Fallout* universe.

---

## The Scenario: Tactical Satellite Sync

You are at the top of a decaying radar dish installation in the Capital Wasteland. You’ve spliced into a **RobCo High-Orbit Tactical Link (HOTL)** terminal. The terminal is receiving a fire-hose stream of real-time telemetry data tracking incoming ballistic threats.

The telemetry packets are being dumped continuously into a fixed-size, circular memory arena (a Ring Buffer). Because the hardware dates back to 2077, it uses a **multithreaded lock-free execution architecture**.

The system writer implemented a custom sliding logging engine to calculate a running metric, but they introduced a subtle performance bottleneck. The system is dropping telemetry packets because the current log analysis algorithm runs in quadratic time ($O(N^2)$). If the buffer fills up before processing is done, the terminal overflows, drops connections, and crashes.

Your mission is to rewrite the kernel parsing routine to operate in absolute **Linear Time ($O(N)$)** and **Constant Auxiliary Space ($O(1)$)**, using low-level bit-twiddling.

---

## The System Architecture: High-Low Telemetry Packets

The data structure is a packed array of 64-bit unsigned integers (`[]uint64`), where each element represents a single satellite telemetry snapshot.

To save register space on the ancient processor, each `uint64` is packed with multiple data pieces using strict bit positions:

* **Bits 0–31 (32 bits):** **Timestamp** (Unix epoch delta).
* **Bits 32–47 (16 bits):** **Signal Strength (RSSI)**. A standard unsigned integer.
* **Bits 48–63 (16 bits):** **Radiation Interference Level (RadCount)**. A standard unsigned integer.

### The Twist: The Faulty Sensor

The satellite's primary radiation sensor is damaged and constantly fluctuates. Instead of reporting the true radiation count, it frequently flips a single bit or drops to zero.

To calculate a clean baseline, the system looks at a **Sliding Window of $K$ elements**. For safety reasons, the military targeting computer requires you to report the **maximum radiation value inside every continuous window of size $K$**.

---

## The Mathematical Constraint & Invariants

You are given the array of packed `uint64` telemetry values and the window size $K$.

1. **Strict $O(N)$ Time Complexity:** You cannot use a nested loop to check all $K$ elements for every window step. You cannot use a standard priority queue/heap (which would result in an $O(N \log K)$ runtime). The entire array must be processed in a single, fluid forward pass where each element is effectively evaluated in $O(1)$ amortized time.
2. **Strict $O(1)$ Auxiliary Space:** **Here is the hardware twist.** In a standard software interview, you would solve this using a Monotonic Deque (Double-Ended Queue) storing indices. However, the RobCo kernel memory allocator is entirely offline. **You cannot allocate a slice, array, or dynamic list to act as your queue.**
* *How do you maintain a sliding maximum across a window of size $K$ without allocating an auxiliary data structure to hold the history?*
* *Hint:* The underlying ring-buffer array is *already* in memory. Can you use a two-pass block-partitioning trick (like the **Segregated Block / Disjoint Interval Max** algorithm) using only a few fixed pointer variables?


3. **Zero Heap Allocations:** No structural allocations allowed. You must return an array of the results, but the memory for that output array is pre-allocated and passed into your function as a pointer slice to protect the stack.

---

## The Low-Level Kernel Blueprint

Implement the high-performance pipeline in Go matching this kernel signature:

```go
package main

import "fmt"

// AnalyzeTelemetry extracts the 16-bit RadCount from the packed uint64 array 
// and populates the out slice with the maximum RadCount for every sliding window of size K.
// Constraints: Time: O(N), Auxiliary Space: O(1) (excluding the provided out slice).
func AnalyzeTelemetry(telemetry []uint64, k int, out []uint16) {
    n := len(telemetry)
    if n == 0 || k <= 0 || len(out) < n-k+1 {
        return
    }

    // Your zero-allocation bit-manipulation and O(N) sliding max logic here
}

```

---

## Challenge Dataset

Test your architecture logic against this telemetry dump:

```go
func main() {
    // Window size K = 3
    k := 3
    
    // Packed uint64 data stream. 
    // Remember: RadCount is stored in the most significant 16 bits (Bits 48-63).
    // Shift right by 48 to extract it!
    stream := []uint64{
        (uint64(10) << 48) | 0x0000FFFFFFFF, // RadCount = 10
        (uint64(45) << 48) | 0x0000FFFFFFFF, // RadCount = 45
        (uint64(20) << 48) | 0x0000FFFFFFFF, // RadCount = 20
        (uint64(35) << 48) | 0x0000FFFFFFFF, // RadCount = 35
        (uint64(30) << 48) | 0x0000FFFFFFFF, // RadCount = 30
        (uint64(50) << 48) | 0x0000FFFFFFFF, // RadCount = 50
        (uint64(15) << 48) | 0x0000FFFFFFFF, // RadCount = 15
    }

    // Pre-allocated output space (N - K + 1 elements)
    out := make([]uint16, len(stream)-k+1)

    AnalyzeTelemetry(stream, k, out)

    // Expected Output: [45, 45, 35, 50, 50]
    fmt.Printf("Parsed Target Trajectories: %v\n", out)
}

```

### Questions to Crack Before Coding:

1. How do you cleanly isolate bits 48–63 using a bitwise mask or shift without touching the lower telemetry data?
2. How do you implement the "Block Maxima" strategy (pre-calculating prefix and suffix maximums within fixed blocks of size $K$) directly within the existing loops without allocating extra memory pools?

This one forces you to marry bit shifting, array partitioning, and sliding pointer boundaries perfectly. Let’s see how you optimize this RobCo terminal!