package planner

// Node represents a node in the graph of possible actions. Each node contains the
// State of the node at that time, and the Action that the node performs.
type Node struct {
	// state is the state of the world at the time of the Node
	state *State
	// gCost is the total cost from the start to this node
	gCost float64
	// hCost is the heeristic cost from the node to the goal
	hCost float64
	// parent is the node that led to this node
	parent *Node
	// action is the Action that the node represents
	action Action
}

// fCost is the total cost of selecting that node
func (n *Node) fCost() float64 {
	return n.gCost + n.hCost
}

// PriorityQueue is a heap used to select the node with the lowest cost
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].fCost() < pq[j].fCost()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}
