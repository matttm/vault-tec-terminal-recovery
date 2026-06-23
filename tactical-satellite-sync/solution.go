package tacticalsat

import "fmt"

// AnalyzeTelemetry extracts the 16-bit RadCount from each packed uint64 telemetry value
// and writes the maximum RadCount for every sliding window of size k into out.
func AnalyzeTelemetry(telemetry []uint64, k int, out []uint16) {
	n := len(telemetry)
	if n == 0 || k <= 0 || k > n {
		return
	}
	pre := make([]uint16, n)
	suf := make([]uint16, n)
	radCounts := make([]uint16, n)
	runningMax := uint16(0)
	for i := 0; i < n; i++ {
		var radCnt uint16 = uint16(telemetry[i] >> 48)
		radCounts[i] = radCnt
		if i%k == 0 {
			runningMax = radCnt
		} else {
			runningMax = max(radCnt, runningMax)
		}
		pre[i] = runningMax
	}
	for i := n - 1; i >= 0; i-- {
		var radCnt uint16 = uint16(telemetry[i] >> 48)
		if i == n-1 || (i+1)%k == 0 {
			runningMax = radCnt
		} else {
			runningMax = max(radCnt, runningMax)
		}
		suf[i] = runningMax
	}
	fmt.Printf("==============================\n")
	fmt.Printf(" rad=%v\n pre=%v\n suf=%v\n", radCounts, pre, suf)
	i, j := 0, k-1
	for k := range out {
		out[k] = max(suf[i], pre[j])
		i++
		j++
	}
}
