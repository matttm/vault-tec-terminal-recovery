package memorycore

// IsEnclaveCoreSafe returns true if the memory core is safe to boot because it
// contains no directed cycles. It returns false when a cycle is detected.
func IsEnclaveCoreSafe(numBlocks int, links [][]int) bool {
	indegree := make([]int, numBlocks)
	adj := make([][]int, numBlocks)

	for _, link := range links {
		from, to := link[0], link[1]
		adj[from] = append(adj[from], to)
		indegree[to]++
	}

	queue := make([]int, 0, numBlocks)
	for node, degree := range indegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	processed := 0
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		processed++

		for _, neighbor := range adj[cur] {
			indegree[neighbor]--
			if indegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	return processed == numBlocks
}
