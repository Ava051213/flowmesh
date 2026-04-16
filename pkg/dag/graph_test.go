package dag

import "testing"

func TestNewGraph_InitializesNodesMap(t *testing.T) {
	g := NewGraph()
	if g.Nodes == nil {
		t.Fatal("Nodes map should not be nil")
	}
}

func TestGraph_ValidateDepsExist_OK(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A"})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})

	if err := g.ValidateDepsExist(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGraph_ValidateDepsExist_MissingDep(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})

	err := g.ValidateDepsExist()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*MissingDependencyError); !ok {
		t.Fatalf("expected *MissingDependencyError, got %T: %v", err, err)
	}
}

func TestGraph_addNodeStrict_EmptyID(t *testing.T) {
	g := NewGraph()
	if err := g.AddNodeStrict(&Node{ID: ""}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGraph_AddNodeStrict_Duplicate(t *testing.T) {
	g := NewGraph()
	if err := g.AddNodeStrict(&Node{ID: "A"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := g.AddNodeStrict(&Node{ID: "A"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGraph_AddNodeStrict_OK(t *testing.T) {
	g := NewGraph()
	if err := g.AddNodeStrict(&Node{ID: "A"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGraph_ValidateAcyclic_OK(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A"})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})
	g.AddNode(&Node{ID: "C", DependsOn: []string{"A"}})

	if err := g.ValidateAcyclic(); err != nil {
		t.Fatalf("expected no error, got %T: %v", err, err)
	}
}

func TestGraph_ValidateAcyclic_Cycle(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A", DependsOn: []string{"B"}})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"C"}})
	g.AddNode(&Node{ID: "C", DependsOn: []string{"A"}})

	err := g.ValidateAcyclic()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*CycleDetectedError); !ok {
		t.Fatalf("expected *CycleDetectedError, got %T: %v", err, err)
	}
}

func TestGraph_ValidateAcyclic_SelfDependency(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A", DependsOn: []string{"A"}})

	err := g.ValidateAcyclic()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*CycleDetectedError); !ok {
		t.Fatalf("expected *CycleDetectedError, got %T: %v", err, err)
	}
}

func TestGraph_Validate_OK(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A"})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})

	if err := g.Validate(); err != nil {
		t.Fatalf("expected no error, got %T: %v", err, err)
	}
}

func TestGraph_Validate_MissingDep(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})

	err := g.Validate()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*MissingDependencyError); !ok {
		t.Fatalf("expected *MissingDependencyError, got %T: %v", err, err)
	}
}

func TestGraph_Validate_Cycle(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A", DependsOn: []string{"B"}})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"C"}})
	g.AddNode(&Node{ID: "C", DependsOn: []string{"A"}})

	err := g.Validate()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*CycleDetectedError); !ok {
		t.Fatalf("expected *CycleDetectedError, got %T: %v", err, err)
	}
}

func TestGraph_TopoSort_LinearChain(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A"})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})
	g.AddNode(&Node{ID: "C", DependsOn: []string{"B"}})

	order, err := g.TopoSort()
	if err != nil {
		t.Fatalf("expected no error, got %T: %v", err, err)
	}
	if len(order) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %v", len(order), order)
	}

	pos := map[string]int{}
	for i, id := range order {
		pos[id] = i
	}

	if pos["A"] >= pos["B"] || pos["B"] >= pos["C"] {
		t.Fatalf("expected A before B before C, got %v", order)
	}
}

func TestGraph_TopoSort_Branching(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A"})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})
	g.AddNode(&Node{ID: "C", DependsOn: []string{"A"}})
	g.AddNode(&Node{ID: "D", DependsOn: []string{"B", "C"}})

	order, err := g.TopoSort()
	if err != nil {
		t.Fatalf("expected no error, got %T: %v", err, err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 nodes, got %d: %v", len(order), order)
	}

	pos := map[string]int{}
	for i, id := range order {
		pos[id] = i
	}
	
	if pos["A"] >= pos["B"] {
		t.Fatalf("expected A before B, got %v", order)
	}
	if pos["A"] >= pos["C"] {
		t.Fatalf("expected A before C, got %v", order)
	}
	if pos["B"] >= pos["D"] || pos["C"] >= pos["D"] {
		t.Fatalf("expected B and C before D, got %v", order)
	}
}

func TestGraph_TopoSort_Cycle(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "A", DependsOn: []string{"B"}})
	g.AddNode(&Node{ID: "B", DependsOn: []string{"A"}})
	
	_, err := g.TopoSort()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*CycleDetectedError); !ok {
		t.Fatalf("expected *CycleDetectedError, got %T: %v", err, err)
	}
}







