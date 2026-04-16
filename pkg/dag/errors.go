package dag

import (
	"fmt"
	"strings"
)

type MissingDependencyError struct {
	NodeID string
	DepID  string
}

func (e *MissingDependencyError) Error() string {
	return fmt.Sprintf("node %s depends on missing node %s", e.NodeID, e.DepID)
}

type InvalidNodeIDError struct{}

func (e *InvalidNodeIDError) Error() string {
	return "node id is empty"
}

type DuplicateNodeError struct {
	NodeID string
}

func (e *DuplicateNodeError) Error() string {
	return fmt.Sprintf("duplicate node id: %s", e.NodeID)
}

type CycleDetectedError struct {
	Cycle []string
}

func (e *CycleDetectedError) Error() string {
	return fmt.Sprintf("cycle detected: %s", strings.Join(e.Cycle, " -> "))
}

