package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"dayplanner/internal/domain"
)

const dataFile = "data.json"

type JSONRepository struct {
	dir string
}

func NewJSONRepository(dir string) (*JSONRepository, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("repository: create config dir %q: %w", dir, err)
	}
	return &JSONRepository{dir: dir}, nil
}

func (r *JSONRepository) dataPath() string {
	return filepath.Join(r.dir, dataFile)
}

func (r *JSONRepository) tempPath() string {
	return filepath.Join(r.dir, "."+dataFile+".tmp")
}

func (r *JSONRepository) Load() ([]domain.Task, error) {
	data, err := os.ReadFile(r.dataPath())
	if os.IsNotExist(err) {
		return []domain.Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: read %q: %w", r.dataPath(), err)
	}

	var tasks []domain.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("repository: parse %q: %w", r.dataPath(), err)
	}

	return tasks, nil
}

func (r *JSONRepository) Save(tasks []domain.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("repository: marshal tasks: %w", err)
	}

	tmp := r.tempPath()
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("repository: open tmp file %q: %w", tmp, err)
	}

	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return fmt.Errorf("repository: write tmp file %q: %w", tmp, err)
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("repository: sync tmp file %q: %w", tmp, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("repository: close tmp file %q: %w", tmp, err)
	}

	if err := os.Rename(tmp, r.dataPath()); err != nil {
		return fmt.Errorf("repository: rename %q → %q: %w", tmp, r.dataPath(), err)
	}

	return nil
}
