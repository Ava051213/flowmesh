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
