package pathfinding

import "container/heap"

type Path struct {
	CurrentPos int
	Nodes      []*AStarNode
}

type AStar struct {
	start, end *AStarNode
}

func (a *AStar) FindPath(start, end *git
) Path {

	openSet := AStarHeap{start}
	cameFrom := make(map[*AStarNode]*AStarNode)
	gScore := make(map[*AStarNode]int)
	gScore[start] = 0

	fScore := make(map[*AStarNode]int)
	fScore[start] = a.heuristic(start)

	for len(openSet) > 0 {
		current := heap.Pop(&openSet).(*AStarNode)
		if current == end {
			return Path{} // TODO: return actual path
		}
		for neighbor := range current.GetNeighbors() {

		}
	}

	return Path{}
}

func (a *AStar) heuristic(*AStarNode) int {
	return 0
}

type AStarNode struct {
	hScore int
	gScore int
}

func (a *AStarNode) GetNeighbors() []AStarNode {
	return []AStarNode{}
}

func (a *AStarNode) fScore() int {
	return a.hScore + a.gScore
}

type AStarHeap []*AStarNode

func (h AStarHeap) Len() int           { return len(h) }
func (h AStarHeap) Less(i, j int) bool { return h[i].fScore() > h[j].fScore() }
func (h AStarHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *AStarHeap) Push(x any) {
	*h = append(*h, x.(*AStarNode))
}

func (h *AStarHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
