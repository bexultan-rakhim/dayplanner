package scheduler

import (
	"testing"
	"time"

	"dayplanner/internal/domain"
	"dayplanner/internal/graph"
)

func makeTask(id string, priority domain.Priority, createdAt time.Time, deps ...string) domain.Task {
	return domain.Task{
		ID:        id,
		Tag:       "TST",
		Name:      id,
		Priority:  priority,
		Status:    domain.StatusTodo,
		DependsOn: deps,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

var (
	t0 = time.Date(2025, 1, 1, 8, 0, 0, 0, time.UTC)
	t1 = t0.Add(time.Minute)
	t2 = t0.Add(2 * time.Minute)
)

func mustBuild(t *testing.T, tasks []domain.Task) *graph.Graph {
	t.Helper()
	g, err := graph.Build(tasks)
	if err != nil {
		t.Fatalf("graph.Build: %v", err)
	}
	return g
}

func ids(tasks []domain.Task) []string {
	out := make([]string, len(tasks))
	for i, t := range tasks {
		out[i] = t.ID
	}
	return out
}

func TestOrder_Empty(t *testing.T) {
	g := mustBuild(t, nil)
	if result := Order(nil, g); result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestOrder_SingleTask(t *testing.T) {
	tasks := []domain.Task{makeTask("TST-a", domain.PriorityHigh, t0)}
	g := mustBuild(t, tasks)
	result := Order(tasks, g)
	if len(result) != 1 || result[0].ID != "TST-a" {
		t.Errorf("unexpected result: %v", ids(result))
	}
}

func TestOrder_LayerBeforePriority(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-a", domain.PriorityLow, t0),
		makeTask("TST-b", domain.PriorityHigh, t1, "TST-a"),
	}
	g := mustBuild(t, tasks)
	result := Order(tasks, g)
	if result[0].ID != "TST-a" || result[1].ID != "TST-b" {
		t.Errorf("dependency layer should take precedence over priority, got %v", ids(result))
	}
}

func TestOrder_PriorityWithinLayer(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-low", domain.PriorityLow, t0),
		makeTask("TST-high", domain.PriorityHigh, t1),
		makeTask("TST-med", domain.PriorityMedium, t2),
	}
	g := mustBuild(t, tasks)
	result := Order(tasks, g)
	expected := []string{"TST-high", "TST-med", "TST-low"}
	for i, id := range expected {
		if result[i].ID != id {
			t.Errorf("position %d: want %q, got %q", i, id, result[i].ID)
		}
	}
}

func TestOrder_CreatedAtTiebreaker(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-later", domain.PriorityHigh, t2),
		makeTask("TST-earlier", domain.PriorityHigh, t0),
		makeTask("TST-middle", domain.PriorityHigh, t1),
	}
	g := mustBuild(t, tasks)
	result := Order(tasks, g)
	expected := []string{"TST-earlier", "TST-middle", "TST-later"}
	for i, id := range expected {
		if result[i].ID != id {
			t.Errorf("position %d: want %q, got %q", i, id, result[i].ID)
		}
	}
}

func TestOrder_DiamondDependency(t *testing.T) {
	//   TST-a (layer 0)
	//   /         \
	// TST-b       TST-c  (layer 1)
	//   \         /
	//    TST-d (layer 2)
	tasks := []domain.Task{
		makeTask("TST-a", domain.PriorityHigh, t0),
		makeTask("TST-b", domain.PriorityHigh, t1, "TST-a"),
		makeTask("TST-c", domain.PriorityHigh, t2, "TST-a"),
		makeTask("TST-d", domain.PriorityHigh, t0, "TST-b", "TST-c"),
	}
	g := mustBuild(t, tasks)
	result := Order(tasks, g)

	layerOf := func(id string) int {
		for i, task := range result {
			if task.ID == id {
				return i
			}
		}
		return -1
	}

	if layerOf("TST-a") >= layerOf("TST-b") {
		t.Error("TST-a must come before TST-b")
	}
	if layerOf("TST-a") >= layerOf("TST-c") {
		t.Error("TST-a must come before TST-c")
	}
	if layerOf("TST-b") >= layerOf("TST-d") {
		t.Error("TST-b must come before TST-d")
	}
	if layerOf("TST-c") >= layerOf("TST-d") {
		t.Error("TST-c must come before TST-d")
	}
}

func TestOrder_DoesNotMutateInput(t *testing.T) {
	tasks := []domain.Task{
		makeTask("TST-low", domain.PriorityLow, t0),
		makeTask("TST-high", domain.PriorityHigh, t1),
	}
	original := make([]domain.Task, len(tasks))
	copy(original, tasks)

	g := mustBuild(t, tasks)
	Order(tasks, g)

	for i := range tasks {
		if tasks[i].ID != original[i].ID {
			t.Errorf("input slice was mutated at index %d", i)
		}
	}
}
