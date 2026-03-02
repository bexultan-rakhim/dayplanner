package repository

import "dayplanner/internal/domain"

type Repository interface {
	Load() ([]domain.Task, error)
	Save(tasks []domain.Task) error
}

