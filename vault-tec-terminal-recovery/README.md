Here is a multi-layered, low-level (LV) systems programming puzzle wrapped completely in the lore of the *Fallout* universe.

To solve this, you can't just think like a high-level developer; you have to think like a RobCo hardware engineer trying to salvage corrupted memory banks from a terminal inside an abandoned Vault.

---

## The Scenario: The Vault-Tec Terminal Recovery

You are exploring the dark, radioactive ruins of **Vault 97**. You find a functional but flickering RobCo Industries terminal connected to the main vault door override mechanism.

The security system has intercepted a manual override command, but the terminal's RAM was heavily blasted by electromagnetic pulse (EMP) radiation when the bombs dropped. The data packet containing the override passcode is trapped inside a raw byte array.

The automated terminal interface is locked, displaying only a diagnostic screen. It expects you to write a high-throughput, zero-allocation decoding routine to parse the transmission stream. If you allocate a single heap object, the terminal's primitive 640KB operating system will crash, locking the Vault doors forever.

---

## The Protocol Specification: Vault-Net v2.1

The raw data stream arrives as a continuous slice of bytes (`[]byte`). The protocol uses a tight, variable-length packet encoding scheme to conserve precious copper wiring.

A valid override command structure consists of three sequential fields:

### 1. The Preamble (4 Bytes)

Every valid Vault-Tec transmission must begin with the magic hex signature representing the classic RobCo handshake:

* `0x52 0x4F 0x42 0x43` (ASCII for `ROBC`)

### 2. The Metadata Bitmask (2 Bytes / 16-bit Unsigned Integer)

The next two bytes form a **Big-Endian `uint16**` bitmask containing system flags. You must parse this integer and extract three specific packed properties using bitwise operations:

* **Bits 0–3 (4 bits):** The Clearance Level. (Must be exactly `0x0F` / Overseer Clearance).
* **Bits 4–11 (8 bits):** The Sector ID. (Must be exactly `0x65` / Reactor Sector).
* **Bits 12–15 (4 bits):** Encryption Type.
* If it is `0x0`, the payload is plaintext.
* If it is `0x1`, every byte in the subsequent payload payload has been obfuscated via a bitwise **XOR operation with a constant key of `0x55**` (Rad-Away standard).



### 3. The Payload Length (1 Byte)

A single byte indicating the exact length ($N$) of the string passcode that immediately follows.

### 4. The Passcode Payload ($N$ Bytes)

The remaining $N$ bytes represent the ASCII characters of the passcode string.

---

## The Puzzle Constraints

To prove your low-level engineering chops, you must implement a Go function that successfully extracts the passcode under strict performance invariants:

1. **Zero Allocations ($O(1)$ Space):** You cannot create new slices, maps, or objects. You must parse the existing data in-place.
2. **Zero Inefficient Casting:** You cannot cast the byte array into a string until the final validation step. Everything must be evaluated at the raw byte level.
3. **Robust Boundary Defenses:** The data stream is highly unstable. If the payload length byte claims the passcode is 20 bytes long, but the slice only has 5 bytes remaining, your code must safely return an error without panicking (no index out of bounds).

---

## The Input Data Stream

Here is the exact corrupted raw byte slice recovered from the terminal's memory buffer:

```go
stream := []byte{
    0x52, 0x4F, 0x42, 0x43, // Preamble
    0x16, 0x5F,             // Metadata Bitmask
    0x0C,                   // Payload Length
    0x13, 0x34, 0x37, 0x37, // Encrypted Payload Bytes...
    0x3A, 0x20, 0x13, 0x3A, 
    0x3F, 0x39, 0x33, 0x23,
}

```

---

## Your Mission

Write the decoding function in Go matching this signature:

```go
// DecodeVaultOverride parses the raw stream in-place.
// It returns the decoded passcode string, or an error if the packet is corrupt or unauthorized.
func DecodeVaultOverride(stream []byte) (string, error) {
    // Your low-level logic here
}

```

### Questions to Answer:

1. Is the data stream provided above actually valid according to the Overseer Clearance and Reactor Sector bitmask requirements?
2. What is the encryption state of this specific packet?
3. If you run your bitwise arithmetic and decryption logic over the stream, **what is the hidden ASCII passcode string that opens the Vault?**

*Grab your Pip-Boy, fire up the compiler, and show me the raw byte manipulation!*

---

## Local LeetCode-Style Setup

This repo is set up as a small Go coding challenge. The entry point is:

```go
func DecodeVaultOverride(stream []byte) (string, error)
```

Use these commands while working:

```sh
go test ./...
go run ./cmd/decode
```

Files:

* `solution.go` contains the starter stub for your implementation.
* `solution_test.go` contains visible test cases for valid packets, corrupt packets, and zero-allocation successful decoding. These tests are expected to fail until you implement the stub.
* `cmd/decode/main.go` runs the recovered byte stream from this README and prints the decoded passcode.
