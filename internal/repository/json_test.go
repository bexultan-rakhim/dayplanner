package repository

import (
	"os"
	"testing"
	"time"

	"dayplanner/internal/domain"
)

func newTestRepo(t *testing.T) *JSONRepository {
	t.Helper()
	dir := t.TempDir()
	repo, err := NewJSONRepository(dir)
	if err != nil {
		t.Fatalf("NewJSONRepository: %v", err)
	}
	return repo
}

func sampleTask(id, tag string) domain.Task {
	now := time.Now().UTC().Truncate(time.Second)
	return domain.Task{
		ID:        id,
		Tag:       tag,
		Name:      "Sample task",
		Goal:      "It works",
		Priority:  domain.PriorityHigh,
		Status:    domain.StatusTodo,
		DependsOn: []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestLoad_NoFileReturnsEmpty(t *testing.T) {
	repo := newTestRepo(t)
	tasks, err := repo.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty slice, got %d tasks", len(tasks))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	repo := newTestRepo(t)
	original := []domain.Task{
		sampleTask("AUTH-login", "AUTH"),
		sampleTask("API-schema", "API"),
	}

	if err := repo.Save(original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(original) {
		t.Fatalf("expected %d tasks, got %d", len(original), len(loaded))
	}
	for i := range original {
		if loaded[i].ID != original[i].ID {
			t.Errorf("[%d] ID: want %q, got %q", i, original[i].ID, loaded[i].ID)
		}
		if loaded[i].Status != original[i].Status {
			t.Errorf("[%d] Status: want %q, got %q", i, original[i].Status, loaded[i].Status)
		}
		if loaded[i].Priority != original[i].Priority {
			t.Errorf("[%d] Priority: want %q, got %q", i, original[i].Priority, loaded[i].Priority)
		}
		if !loaded[i].CreatedAt.Equal(original[i].CreatedAt) {
			t.Errorf("[%d] CreatedAt: want %v, got %v", i, original[i].CreatedAt, loaded[i].CreatedAt)
		}
	}
}

func TestSaveAndLoad_DependsOnPreserved(t *testing.T) {
	repo := newTestRepo(t)
	task := sampleTask("AUTH-deploy", "AUTH")
	task.DependsOn = []string{"AUTH-login", "AUTH-tests"}

	if err := repo.Save([]domain.Task{task}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded[0].DependsOn) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(loaded[0].DependsOn))
	}
	if loaded[0].DependsOn[0] != "AUTH-login" || loaded[0].DependsOn[1] != "AUTH-tests" {
		t.Errorf("unexpected DependsOn: %v", loaded[0].DependsOn)
	}
}

func TestSave_OverwritesPreviousData(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Save([]domain.Task{sampleTask("AUTH-login", "AUTH")}); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	if err := repo.Save([]domain.Task{sampleTask("API-schema", "API")}); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 1 || loaded[0].ID != "API-schema" {
		t.Errorf("expected only API-schema, got %v", loaded)
	}
}

func TestSave_EmptySlice(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Save([]domain.Task{sampleTask("AUTH-login", "AUTH")}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := repo.Save([]domain.Task{}); err != nil {
		t.Fatalf("Save empty: %v", err)
	}

	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("expected empty slice after saving empty, got %d tasks", len(loaded))
	}
}

func TestLoad_CorruptFileReturnsError(t *testing.T) {
	repo := newTestRepo(t)

	if err := os.WriteFile(repo.dataPath(), []byte("not json {{{"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := repo.Load()
	if err == nil {
		t.Error("expected error on corrupt JSON, got nil")
	}
}

func TestNewJSONRepository_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	nested := dir + "/a/b/c"

	_, err := NewJSONRepository(nested)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(nested); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
