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

		aDone := a.Status == domain.StatusDone
		bDone := b.Status == domain.StatusDone
		if aDone != bDone {
			return !aDone
		}

		if layerOf[a.ID] != layerOf[b.ID] {
			return layerOf[a.ID] < layerOf[b.ID]
		}
		if a.Priority != b.Priority {
			return priorityRank(a.Priority) < priorityRank(b.Priority)
		}
		return a.CreatedAt.Before(b.CreatedAt)
	})

	return ordered
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
