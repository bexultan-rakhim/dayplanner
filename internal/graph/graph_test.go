package graph

import (
	"testing"
	"time"

	"dayplanner/internal/domain"
)

func makeTask(id string, deps ...string) domain.Task {
	return domain.Task{
		ID:        id,
		Tag:       "TST",
		Name:      id,
		Priority:  domain.PriorityMedium,
		Status:    domain.StatusTodo,
		DependsOn: deps,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func TestBuild_Empty(t *testing.T) {
	g, err := Build([]domain.Task{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if layers := g.Layers(); len(layers) != 0 {
		t.Errorf("expected no layers, got %d", len(layers))
	}
}

func TestBuild_UnknownDependency(t *testing.T) {
	tasks := []domain.Task{makeTask("TST-a", "TST-unknown")}
	_, err := Build(tasks)
	if err == nil {
		t.Error("expected error for unknown dependency")
	}
}

func TestBuild_CycleDetected(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-a", "TST-b"),
		makeTask("TST-b", "TST-a"),
	}
	_, err := Build(tasks)
	if err == nil {
		t.Error("expected error for cycle")
	}
}

func TestBuild_SelfCycle(t *testing.T) {
	tasks := []domain.Task{makeTask("TST-a", "TST-a")}
	_, err := Build(tasks)
	if err == nil {
		t.Error("expected error for self-cycle")
	}
}

func TestLayers_LinearChain(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-a"),
		makeTask("TST-b", "TST-a"),
		makeTask("TST-c", "TST-b"),
	}
	g, err := Build(tasks)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	layers := g.Layers()
	if len(layers) != 3 {
		t.Fatalf("expected 3 layers, got %d", len(layers))
	}
	if layers[0][0] != "TST-a" {
		t.Errorf("layer 0 should be TST-a, got %v", layers[0])
	}
	if layers[1][0] != "TST-b" {
		t.Errorf("layer 1 should be TST-b, got %v", layers[1])
	}
	if layers[2][0] != "TST-c" {
		t.Errorf("layer 2 should be TST-c, got %v", layers[2])
	}
}

func TestLayers_DiamondUsesLongestPath(t *testing.T) {
	//       TST-a
	//      /     \
	//   TST-b   TST-c
	//      \     /
	//       TST-d
	tasks := []domain.Task{
		makeTask("TST-a"),
		makeTask("TST-b", "TST-a"),
		makeTask("TST-c", "TST-a"),
		makeTask("TST-d", "TST-b", "TST-c"),
	}
	g, err := Build(tasks)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	layers := g.Layers()
	if len(layers) != 3 {
		t.Fatalf("expected 3 layers for diamond, got %d: %v", len(layers), layers)
	}
	layerOf := func(id string) int {
		for i, l := range layers {
			for _, n := range l {
				if n == id {
					return i
				}
			}
		}
		return -1
	}
	if layerOf("TST-d") != 2 {
		t.Errorf("TST-d should be in layer 2, got layer %d", layerOf("TST-d"))
	}
}

func TestBlocking_And_BlockedBy(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-a"),
		makeTask("TST-b", "TST-a"),
		makeTask("TST-c", "TST-a"),
	}
	g, err := Build(tasks)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	blockedBy := g.BlockedBy("TST-a")
	if len(blockedBy) != 2 {
		t.Errorf("TST-a should block 2 tasks, got %v", blockedBy)
	}

	blocking := g.Blocking("TST-b")
	if len(blocking) != 1 || blocking[0] != "TST-a" {
		t.Errorf("TST-b should be blocked by TST-a, got %v", blocking)
	}

	if len(g.Blocking("TST-a")) != 0 {
		t.Errorf("TST-a has no dependencies, got %v", g.Blocking("TST-a"))
	}
}

func TestCriticalPath_LinearChain(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-a"),
		makeTask("TST-b", "TST-a"),
		makeTask("TST-c", "TST-b"),
	}
	g, _ := Build(tasks)
	path := g.CriticalPath()
	if len(path) != 3 {
		t.Fatalf("expected critical path length 3, got %d: %v", len(path), path)
	}
	expected := []string{"TST-a", "TST-b", "TST-c"}
	for i, id := range expected {
		if path[i] != id {
			t.Errorf("path[%d]: want %q, got %q", i, id, path[i])
		}
	}
}

func TestCriticalPath_PicksLongest(t *testing.T) {
	// TST-a → TST-b → TST-d  (length 3)
	// TST-c → TST-d           (length 2)
	tasks := []domain.Task{
		makeTask("TST-a"),
		makeTask("TST-b", "TST-a"),
		makeTask("TST-c"),
		makeTask("TST-d", "TST-b", "TST-c"),
	}
	g, _ := Build(tasks)
	path := g.CriticalPath()
	if len(path) != 3 {
		t.Fatalf("expected critical path length 3, got %d: %v", len(path), path)
	}
	if path[0] != "TST-a" || path[len(path)-1] != "TST-d" {
		t.Errorf("unexpected critical path: %v", path)
	}
}

func TestCriticalPath_Empty(t *testing.T) {
	g, _ := Build([]domain.Task{})
	if path := g.CriticalPath(); path != nil {
		t.Errorf("expected nil critical path for empty graph, got %v", path)
	}
}
