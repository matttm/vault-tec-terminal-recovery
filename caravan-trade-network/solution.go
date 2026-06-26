package caravan

import (
	"container/heap"
)

type routeState struct {
	node int
	risk int
}

type riskHeap []routeState

var _ heap.Interface = (*riskHeap)(nil)

func (h riskHeap) Len() int {
	return len(h)
}

func (h riskHeap) Less(i, j int) bool {
	return h[i].risk < h[j].risk
}

func (h riskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *riskHeap) Push(x any) {
	*h = append(*h, x.(routeState))
}

func (h *riskHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// MinimizeMaxCaravanRisk finds a path from source to destination that minimizes
// the maximum risk score of any single route along that path.
func MinimizeMaxCaravanRisk(n int, routes [][]int, source int, destination int) int {
	// Creating adjancency list
	h := riskHeap {}
	adjList := map[int][]routeState {}
	for _, route := range routes {
		src, dst, rsk := route[0], route[1], route[2]
		if len(adjList[src]) == 0 {
			adjList[src] = []routeState{}
		}
		adjList[src] = append(adjList[src], routeState{ node: dst, risk: rsk })
	}
	// Set starting node
	heap.Push(&h, routeState { node: source, risk: 0 })
	for h.Len() > 0 {
		state := heap.Pop(&h).(routeState)
		cur := state.node
		rsk := state.risk
		if cur == destination {
			return rsk
		}
		for _, neighbor := range adjList[cur] {
			heap.Push(&h, routeState { node: neighbor.node, risk: max(neighbor.risk, rsk)})
		}
		adjList[cur] = nil
	}
	return -1
}
