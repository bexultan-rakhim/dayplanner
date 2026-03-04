package scheduler

import (
	"sort"

	"dayplanner/internal/domain"
	"dayplanner/internal/graph"
)

func Order(tasks []domain.Task, g *graph.Graph) []domain.Task {
	if len(tasks) == 0 {
		return nil
	}

	layers := g.Layers()
	layerOf := make(map[string]int, len(tasks))
	for i, layer := range layers {
		for _, id := range layer {
			layerOf[id] = i
		}
	}

	ordered := make([]domain.Task, len(tasks))
	copy(ordered, tasks)

	sort.SliceStable(ordered, func(i, j int) bool {
		a, b := ordered[i], ordered[j]

		if a.Status != b.Status {
			return priorityStatus(a.Status) < priorityStatus(b.Status)
		}

		if a.Priority != b.Priority {
			return priorityRank(a.Priority) < priorityRank(b.Priority)
		}

		if layerOf[a.ID] != layerOf[b.ID] {
			return layerOf[a.ID] < layerOf[b.ID]
		}
		return a.CreatedAt.Before(b.CreatedAt)
	})

	return ordered
}

func priorityStatus(s domain.Status) int {
	switch s {
	case domain.StatusInProgress:
		return 0
	case domain.StatusTodo:
		return 1
	case domain.StatusBlocked:
		return 2
	case domain.StatusDone:
		return 3
	default:
		return 0
	}
}

func priorityRank(p domain.Priority) int {
	switch p {
	case domain.PriorityHigh:
		return 0
	case domain.PriorityMedium:
		return 1
	case domain.PriorityLow:
		return 2
	default:
		return 3
	}
}
