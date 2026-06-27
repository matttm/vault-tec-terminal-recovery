package memorycore

// IsEnclaveCoreSafe returns true if the memory core is safe to boot because it
// contains no directed cycles. It returns false when a cycle is detected.
func IsEnclaveCoreSafe(numBlocks int, links [][]int) bool {
	seen := make([]bool, numBlocks)
	indegree := make([]int, numBlocks)
	outdegree := make([]int, numBlocks)
	adj := make([][]int, numBlocks)
	seenCnt := 0
	queue := []int {}
	// getting root nodes
	for _, link := range links {
		indegree[link[1]]++
		outdegree[link[0]]++
	}
	for i, v := range indegree {
		if v == 0 {
			queue = append(queue, i)
		}
	}
	// construct adjacency list
	for _, link := range links {
		adj[link[0]] = append(adj[link[0]], link[1])
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if seen[cur] {
			return false
		}
		seen[cur] = true
		seenCnt++
		for _, neighbor := range adj[cur] {
			if seen[neighbor] {
				return false
			}
			queue = append(queue, neighbor)
		}
	}
	return seenCnt == numBlocks
}
