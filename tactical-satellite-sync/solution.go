package tacticalsat

// AnalyzeTelemetry extracts the 16-bit RadCount from each packed uint64 telemetry value
// and writes the maximum RadCount for every sliding window of size k into out.
func AnalyzeTelemetry(telemetry []uint64, k int, out []uint16) {
	n := len(telemetry)
	if n == 0 || k <= 0 || k > n || len(out) < n-k+1 {
		return
	}

	windows := n - k + 1
	runningMax := uint16(0)
	for i := 0; i < n; i++ {
		radCnt := uint16(telemetry[i] >> 48)
		if i%k == 0 {
			runningMax = radCnt
		} else if radCnt > runningMax {
			runningMax = radCnt
		}

		// Store only the prefix maximums needed by each result window in the
		// preallocated output array instead of allocating a full prefix slice.
		if i >= k-1 {
			out[i-k+1] = runningMax
		}
	}

	for i := n - 1; i >= 0; i-- {
		radCnt := uint16(telemetry[i] >> 48)
		if i == n-1 || (i+1)%k == 0 {
			runningMax = radCnt
		} else if radCnt > runningMax {
			runningMax = radCnt
		}

		// Fold the matching suffix maximum into the result in-place, completing
		// max(suffix[start], prefix[end]) without allocating a suffix slice.
		if i < windows && runningMax > out[i] {
			out[i] = runningMax
		}
	}
}
