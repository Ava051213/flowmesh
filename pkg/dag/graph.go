package dag

import "fmt"

type Node struct {
	ID        string
	DependsOn []string
	Metadata  map[string]string
}

type Graph struct {
	Nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{Nodes: map[string]*Node{}}
}

func (g *Graph) AddNode(n *Node) {
	if g.Nodes == nil {
		g.Nodes = map[string]*Node{}
	}
	g.Nodes[n.ID] = n
}

func (g *Graph) ValidateDepsExist() error {
	for _, n := range g.Nodes {
		for _, depID := range n.DependsOn {
			if _, ok := g.Nodes[depID]; !ok {
				return &MissingDependencyError{NodeID: n.ID, DepID: depID}
			}
		}
	}
	return nil
}

func (g *Graph) AddNodeStrict(n *Node) error {
	if n == nil || n.ID == "" {
		return &InvalidNodeIDError{}
	}
	if g.Nodes == nil {
		g.Nodes = map[string]*Node{}
	}
	if _, exists := g.Nodes[n.ID]; exists {
		return &DuplicateNodeError{NodeID: n.ID}
	}
	g.Nodes[n.ID] = n
	return nil
}

func (g *Graph) ValidateAcyclic() error {
	// 0=未访问 1=访问中 2=已完成
	state := make(map[string]uint8, len(g.Nodes))

	// 用递归栈构造环路径
	stack := make([]string, 0, len(g.Nodes))
	pos := make(map[string]int, len(g.Nodes)) // 节点在 stack 中的位置

	var visit func(id string) error
	visit = func(id string) error {
		switch state[id] {
		case 2:
			return nil
		case 1:
			start := pos[id]
			cycle := append([]string(nil), stack[start:]...)
			cycle = append(cycle, id) // 首尾同节点，明确闭环
			return &CycleDetectedError{Cycle: cycle}
		}

		n, ok := g.Nodes[id]
		if !ok {
			// 理论上应由 ValidateDepsExist 提前拦住；这里做防御。
			return &MissingDependencyError{NodeID: id, DepID: id}
		}

		state[id] = 1
		pos[id] = len(stack)
		stack = append(stack, id)

		for _, depID := range n.DependsOn {
			if _, ok := g.Nodes[depID]; !ok {
				return &MissingDependencyError{NodeID: id, DepID: depID}
			}
			if err := visit(depID); err != nil {
				return err
			}
		}

		stack = stack[:len(stack)-1]
		delete(pos, id)
		state[id] = 2
		return nil
	}

	for id := range g.Nodes {
		if state[id] != 0 {
			continue
		}
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) Validate() error {
	if g == nil {
		return fmt.Errorf("graph is nil")
	}
	if g.Nodes == nil {
		return nil
	}
	if err := g.ValidateDepsExist(); err != nil {
		return err
	}
	if err := g.ValidateAcyclic(); err != nil {
		return err
	}
	return nil
}

func (g *Graph) TopoSort() ([]string, error) {
	if err := g.Validate(); err != nil { return nil, err }

	if g.Nodes == nil {
		return nil, nil
	}

	indegree := make(map[string]int, len(g.Nodes))
	out := make(map[string][]string, len(g.Nodes))

	for id := range g.Nodes {
		indegree[id] = 0
	}

	for id, n := range g.Nodes {
		for _, depID := range n.DependsOn {
			indegree[id]++
			out[depID] = append(out[depID], id)
		}
	}

	queue := make([]string, 0, len(g.Nodes))
	for id, d := range indegree {
		if d == 0 {
			queue = append(queue, id)
		}
	}

	order := make([]string, 0, len(g.Nodes))
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]

		order = append(order, id)

		for _, to := range out[id] {
			indegree[to]--
			if indegree[to] == 0 {
				queue = append(queue, to)
			}
		}
	}

	if len(order) != len(g.Nodes) {
		return nil, &CycleDetectedError{Cycle: []string{}}
	}
	return order, nil

	return nil, fmt.Errorf("not implemented")
}

