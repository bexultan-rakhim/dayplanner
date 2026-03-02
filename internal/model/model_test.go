package model

import (
	"fmt"
	"testing"
	"time"

	"dayplanner/internal/domain"
	"dayplanner/internal/history"
	"dayplanner/internal/repository"
)

func sampleTask(id, tag string, priority domain.Priority, deps ...string) domain.Task {
	now := time.Now().UTC()
	return domain.Task{
		ID:        id,
		Tag:       tag,
		Name:      id,
		Priority:  priority,
		Status:    domain.StatusTodo,
		DependsOn: deps,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func newTestModel(tasks []domain.Task) Model {
	repo := repository.NewInMemoryRepository()
	_ = repo.Save(tasks)
	m := Model{
		Page:    PageDashboard,
		Tasks:   tasks,
		History: history.New(),
		Repo:    repo,
	}
	return m.RebuildGraph()
}

func TestNew_LoadsTasksFromRepo(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	tasks := []domain.Task{sampleTask("AUTH-a", "AUTH", domain.PriorityHigh)}
	_ = repo.Save(tasks)

	m, err := New(repo)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(m.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(m.Tasks))
	}
}

func TestRebuildGraph_SetsScheduled(t *testing.T) {
	tasks := []domain.Task{
		sampleTask("AUTH-a", "AUTH", domain.PriorityLow),
		sampleTask("AUTH-b", "AUTH", domain.PriorityHigh, "AUTH-a"),
	}
	m := newTestModel(tasks)

	if m.Graph == nil {
		t.Error("expected graph to be built")
	}
	if len(m.Scheduled) != 2 {
		t.Errorf("expected 2 scheduled tasks, got %d", len(m.Scheduled))
	}
	if m.Scheduled[0].ID != "AUTH-a" {
		t.Errorf("AUTH-a must come first due to dependency, got %q", m.Scheduled[0].ID)
	}
}

func TestRebuildGraph_CycleSetError(t *testing.T) {
	m := Model{
		Tasks:   []domain.Task{},
		History: history.New(),
	}
	m.Tasks = []domain.Task{
		{ID: "AUTH-a", Tag: "AUTH", DependsOn: []string{"AUTH-b"}, CreatedAt: time.Now().UTC()},
		{ID: "AUTH-b", Tag: "AUTH", DependsOn: []string{"AUTH-a"}, CreatedAt: time.Now().UTC()},
	}
	result := m.RebuildGraph()
	if result.Err == nil {
		t.Error("expected error for cyclic graph")
	}
}

func TestTaskByID_Found(t *testing.T) {
	m := newTestModel([]domain.Task{sampleTask("AUTH-a", "AUTH", domain.PriorityHigh)})
	task, ok := m.TaskByID("AUTH-a")
	if !ok {
		t.Fatal("expected task to be found")
	}
	if task.ID != "AUTH-a" {
		t.Errorf("unexpected task ID: %q", task.ID)
	}
}

func TestTaskByID_NotFound(t *testing.T) {
	m := newTestModel(nil)
	_, ok := m.TaskByID("AUTH-missing")
	if ok {
		t.Error("expected task not to be found")
	}
}

func TestSelectedTask_ValidCursor(t *testing.T) {
	tasks := []domain.Task{sampleTask("AUTH-a", "AUTH", domain.PriorityHigh)}
	m := newTestModel(tasks)
	m.Cursor = 0

	task, ok := m.SelectedTask()
	if !ok {
		t.Fatal("expected selected task")
	}
	if task.ID != "AUTH-a" {
		t.Errorf("unexpected task: %q", task.ID)
	}
}

func TestSelectedTask_OutOfBounds(t *testing.T) {
	m := newTestModel(nil)
	m.Cursor = 5
	_, ok := m.SelectedTask()
	if ok {
		t.Error("expected no task for out-of-bounds cursor")
	}
}

func TestWithPage_ResetsCursor(t *testing.T) {
	m := newTestModel(nil)
	m.Cursor = 3
	m = m.WithPage(PageGraph)

	if m.Page != PageGraph {
		t.Errorf("expected PageGraph, got %d", m.Page)
	}
	if m.Cursor != 0 {
		t.Errorf("expected cursor reset to 0, got %d", m.Cursor)
	}
}

func TestWithError_And_ClearError(t *testing.T) {
	m := newTestModel(nil)
	m = m.WithError(fmt.Errorf("invalid id %q", "bad"))
	if m.Err == nil {
		t.Error("expected error to be set")
	}
	m = m.ClearError()
	if m.Err != nil {
		t.Error("expected error to be cleared")
	}
}
