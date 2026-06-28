You’ve been assigned to the **Vault-Tec Department of Internal Security and Intelligence**. Your desk is deep within the administrative sub-levels of Vault-Tec Headquarters.

Your mission is to monitor sub-surface communications from Vault Overseers. Some Overseers are starting to catch on to our "Social Preservation Experiments" and are attempting to leak classified telemetry data to outside factions.

To bypass our automatic network firewalls, a rogue Overseer is trying to smuggle stolen files by hiding them inside standard, everyday **Vault-Tec Security Personnel Status Logs**. They are using a low-level technique known as **Steganography**—hiding secret data right inside the unused padding bits of regular system transmissions.

---

## The Scenario: Operation Bit-Siphon

Every hour, the mainframe receives an array of 32-bit unsigned integers (`[]uint32`). Each integer is a status packet representing a single security guard's active equipment configuration, formatted precisely to save bandwidth:

* **Bits 0–7 (8 bits):** Guard Identification Number.
* **Bits 8–15 (8 bits):** Sector Assignment Code.
* **Bits 16–23 (8 bits):** Weapon Status Flags.
* **Bits 24–31 (8 bits):** **The Padding/Reserved Block.**

According to official Vault-Tec documentation, the top 8 bits (Bits 24–31) must *always* be set to zero (`0x00`). However, our intelligence indicates the rogue Overseer is injecting a secret file byte-by-byte into these exact padding positions across consecutive packets.

---

## The Core Constraint & Low-Level Twist

Your task is to write an intelligence compilation routine that strips away the actual guard logs and siphons out those hidden padding bits.

But you have to do it with extreme caution. The Vault-Tec mainframe runs a background telemetry auditor that checks memory bandwidth usage.

### Your Performance Invariants:

1. **Strict $O(N)$ Time Complexity:** You must sweep through the packet stream in a single linear pass.
2. **Absolute Zero Heap Allocations ($O(1)$ Auxiliary Space):** You cannot create a new byte slice or string buffer to hold the extracted data while you are reading it. You must decode and evaluate the hidden bytes **in-place** using the pre-allocated scratchpad slice provided by the supervisor.
3. **The XOR Cipher Trap:** To prevent casual inspection, the Overseer encrypted the siphoned stream. Every extracted byte must be decrypted by performing a bitwise **XOR operation with a sliding mask** that increments with every hidden byte processed.
* *First hidden byte:* XOR with `0x01`
* *Second hidden byte:* XOR with `0x02`
* *Third hidden byte:* XOR with `0x03`... and so on.



---

## The Go Blueprint

Implement the intelligence-gathering utility matching this kernel signature:

```go
package main

import "fmt"

// SiphonIntelStream extracts, decrypts, and writes hidden data from the packet stream
// into the pre-allocated 'decryptedMessage' buffer. 
// It returns the total number of secret bytes safely extracted.
func SiphonIntelStream(packets []uint32, decryptedMessage []byte) int {
    // Your zero-allocation bit shifting and sliding XOR cipher logic here
}

```

---

## Challenge Dataset

Intercept and decode this raw 32-bit transmission packet array:

```go
func main() {
    // Intercepted packet stream from a rogue Vault Overseer
    packets := []uint32{
        (uint32(0x59) << 24) | 0x00ABCDEF, // Packet 0
        (uint32(0x41) << 24) | 0x00123456, // Packet 1
        (uint32(0x4C) << 24) | 0x00778899, // Packet 2
        (uint32(0x56) << 24) | 0x00000000, // Packet 3
        (uint32(0x55) << 24) | 0x00FF00FF, // Packet 4
        (uint32(0x55) << 24) | 0x00AAAAAA, // Packet 5
    }

	// Pre-allocated buffer provided by Intel Command (matches incoming packets size)
	decryptedMessage := make([]byte, len(packets))

	bytesExtracted := SiphonIntelStream(packets, decryptedMessage)

	// Slice down to the actual bytes recovered
	finalPayload := decryptedMessage[:bytesExtracted]

	fmt.Printf("Total Bytes Recovered: %d\n", bytesExtracted)
	fmt.Printf("Decrypted Secret Message: %s\n", string(finalPayload))
	// Expected Output String: "XCORPS"
}

```

---

### Intelligence Briefing Questions:

1. What bitwise shifting operator (`>>` or `<<`) and what specific hexadecimal bitmask do you need to cleanly drop the lower 24 bits of guard data and isolate just the top 8 padding bits?
2. How do you implement the sliding XOR counter so that it correctly mutates step-by-step alongside your loop index without allocating an extra tracking object?

Get to work, Officer. The future of Vault-Tec depends on your discretion.

---

## Local LeetCode-Style Setup

This directory is set up as a small Go coding challenge. The entry point is:

```go
func SiphonIntelStream(packets []uint32, decryptedMessage []byte) int
```

Use these commands from inside `operation-bit-siphon`:

```sh
go test -v ./...
go run ./cmd/siphon
```

Files:

* `solution.go` contains the starter stub for your implementation.
* `solution_test.go` contains visible test cases for the sample dataset, top-byte extraction, short output buffers, empty input, and zero-allocation successful decoding. These tests are expected to fail until you implement the stub.
* `cmd/siphon/main.go` runs the challenge dataset from this README and prints the recovered message.
