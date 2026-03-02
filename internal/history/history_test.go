package history

import (
	"testing"
	"time"

	"dayplanner/internal/domain"
)

func sampleEntry(id string, kind ActionKind) Entry {
	return Entry{
		Kind: kind,
		Snapshot: domain.Task{
			ID:        id,
			Tag:       "TST",
			Name:      "task " + id,
			Priority:  domain.PriorityMedium,
			Status:    domain.StatusTodo,
			DependsOn: []string{},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
}

func TestNew_StartsEmpty(t *testing.T) {
	s := New()
	if s.CanUndo() {
		t.Error("new stack should not have undo entries")
	}
	if s.CanRedo() {
		t.Error("new stack should not have redo entries")
	}
}

func TestUndo_EmptyReturnsNotOk(t *testing.T) {
	s := New()
	_, ok := s.Undo()
	if ok {
		t.Error("Undo on empty stack should return false")
	}
}

func TestRedo_EmptyReturnsNotOk(t *testing.T) {
	s := New()
	_, ok := s.Redo()
	if ok {
		t.Error("Redo on empty stack should return false")
	}
}

func TestPushAndUndo(t *testing.T) {
	s := New()
	e := sampleEntry("TST-a", ActionDelete)
	s.Push(e)

	got, ok := s.Undo()
	if !ok {
		t.Fatal("expected Undo to succeed")
	}
	if got.Snapshot.ID != "TST-a" {
		t.Errorf("expected TST-a, got %q", got.Snapshot.ID)
	}
	if s.CanUndo() {
		t.Error("undo stack should be empty after single undo")
	}
}

func TestUndo_ThenRedo(t *testing.T) {
	s := New()
	s.Push(sampleEntry("TST-a", ActionEdit))

	s.Undo()
	got, ok := s.Redo()
	if !ok {
		t.Fatal("expected Redo to succeed")
	}
	if got.Snapshot.ID != "TST-a" {
		t.Errorf("expected TST-a, got %q", got.Snapshot.ID)
	}
	if s.CanRedo() {
		t.Error("redo stack should be empty after redo")
	}
}

func TestPush_ClearsRedo(t *testing.T) {
	s := New()
	s.Push(sampleEntry("TST-a", ActionEdit))
	s.Undo()

	if !s.CanRedo() {
		t.Fatal("expected redo entry before push")
	}

	s.Push(sampleEntry("TST-b", ActionEdit))
	if s.CanRedo() {
		t.Error("Push should clear redo stack")
	}
}

func TestUndo_LIFO_Order(t *testing.T) {
	s := New()
	s.Push(sampleEntry("TST-a", ActionEdit))
	s.Push(sampleEntry("TST-b", ActionEdit))
	s.Push(sampleEntry("TST-c", ActionEdit))

	for _, expected := range []string{"TST-c", "TST-b", "TST-a"} {
		got, ok := s.Undo()
		if !ok {
			t.Fatalf("expected Undo to succeed for %q", expected)
		}
		if got.Snapshot.ID != expected {
			t.Errorf("expected %q, got %q", expected, got.Snapshot.ID)
		}
	}
}

func TestPush_EvictsOldestAtMaxDepth(t *testing.T) {
	s := New()
	for i := 0; i < maxDepth; i++ {
		s.Push(sampleEntry("TST-a", ActionEdit))
	}
	s.Push(sampleEntry("TST-overflow", ActionEdit))

	if s.UndoDepth() != maxDepth {
		t.Errorf("expected undo depth %d, got %d", maxDepth, s.UndoDepth())
	}
}

func TestDepth_Counts(t *testing.T) {
	s := New()
	s.Push(sampleEntry("TST-a", ActionEdit))
	s.Push(sampleEntry("TST-b", ActionEdit))

	if s.UndoDepth() != 2 {
		t.Errorf("expected undo depth 2, got %d", s.UndoDepth())
	}
	s.Undo()
	if s.RedoDepth() != 1 {
		t.Errorf("expected redo depth 1, got %d", s.RedoDepth())
	}
}

func TestRestoreIndex_ClampsToTaskCount(t *testing.T) {
	e := sampleEntry("TST-a", ActionDelete)
	e.Index = 10

	if got := RestoreIndex(e, 3); got != 3 {
		t.Errorf("expected clamped index 3, got %d", got)
	}
	if got := RestoreIndex(e, 15); got != 10 {
		t.Errorf("expected original index 10, got %d", got)
	}

	if ActionDelete.String() != "delete" {
		t.Errorf("unexpected: %q", ActionDelete.String())
	}
	if ActionEdit.String() != "edit" {
		t.Errorf("unexpected: %q", ActionEdit.String())
	}
	if ActionKind(99).String() != "unknown" {
		t.Errorf("unexpected: %q", ActionKind(99).String())
	}
}
