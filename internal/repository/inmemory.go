package repository

import "dayplanner/internal/domain"

// Compile-time assertions: both implementations must satisfy Repository.
// If either type drifts out of compliance, the build fails immediately with a
// clear error — no test run required to catch the regression.
var _ Repository = (*JSONRepository)(nil)
var _ Repository = (*InMemoryRepository)(nil)

// InMemoryRepository is a Repository implementation that stores tasks in a
// plain Go slice. It is not safe for concurrent use and is intended solely
// for use in tests of packages that depend on the Repository interface
// (e.g. the model and update layers), allowing them to avoid any filesystem I/O.
type InMemoryRepository struct {
	tasks []domain.Task
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{tasks: []domain.Task{}}
}

func (r *InMemoryRepository) Load() ([]domain.Task, error) {
	cp := make([]domain.Task, len(r.tasks))
	copy(cp, r.tasks)
	return cp, nil
}

func (r *InMemoryRepository) Save(tasks []domain.Task) error {
	cp := make([]domain.Task, len(tasks))
	copy(cp, tasks)
	r.tasks = cp
	return nil
}
