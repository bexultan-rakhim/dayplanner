package graph

import (
	"fmt"

	"dayplanner/internal/domain"
)

type Graph struct {
	tasks    map[string]domain.Task
	children map[string][]string
	parents  map[string][]string
}

func Build(tasks []domain.Task) (*Graph, error) {
	g := &Graph{
		tasks:    make(map[string]domain.Task, len(tasks)),
		children: make(map[string][]string, len(tasks)),
		parents:  make(map[string][]string, len(tasks)),
	}
	for _, t := range tasks {
		g.tasks[t.ID] = t
		if g.children[t.ID] == nil {
			g.children[t.ID] = []string{}
		}
		if g.parents[t.ID] == nil {
			g.parents[t.ID] = []string{}
		}
	}
	for _, t := range tasks {
		for _, dep := range t.DependsOn {
			if _, ok := g.tasks[dep]; !ok {
				return nil, fmt.Errorf("graph: task %q depends on unknown task %q", t.ID, dep)
			}
			g.children[dep] = append(g.children[dep], t.ID)
			g.parents[t.ID] = append(g.parents[t.ID], dep)
		}
	}
	if cycle := g.findCycle(); cycle != "" {
		return nil, fmt.Errorf("graph: cycle detected involving %q", cycle)
	}
	return g, nil
}

func (g *Graph) Layers() [][]string {
	if len(g.tasks) == 0 {
		return nil
	}
	depth := make(map[string]int, len(g.tasks))
	for id := range g.tasks {
		depth[id] = g.longestPath(id, map[string]int{})
	}
	maxDepth := 0
	for _, d := range depth {
		if d > maxDepth {
			maxDepth = d
		}
	}
	layers := make([][]string, maxDepth+1)
	for id, d := range depth {
		layers[d] = append(layers[d], id)
	}
	return layers
}

func (g *Graph) Blocking(id string) []string {
	return g.parents[id]
}

func (g *Graph) BlockedBy(id string) []string {
	return g.children[id]
}

func (g *Graph) CriticalPath() []string {
	terminal := g.terminalNodes()
	if len(terminal) == 0 {
		return nil
	}

	bestLen := -1
	var bestPath []string
	for _, t := range terminal {
		path := g.longestPathTo(t)
		if len(path) > bestLen {
			bestLen = len(path)
			bestPath = path
		}
	}
	return bestPath
}

func (g *Graph) longestPath(id string, memo map[string]int) int {
	if v, ok := memo[id]; ok {
		return v
	}
	max := 0
	for _, dep := range g.parents[id] {
		if d := g.longestPath(dep, memo) + 1; d > max {
			max = d
		}
	}
	memo[id] = max
	return max
}

func (g *Graph) longestPathTo(id string) []string {
	if len(g.parents[id]) == 0 {
		return []string{id}
	}
	var best []string
	for _, dep := range g.parents[id] {
		path := g.longestPathTo(dep)
		if len(path) > len(best) {
			best = path
		}
	}
	return append(best, id)
}

func (g *Graph) terminalNodes() []string {
	var out []string
	for id := range g.tasks {
		if len(g.children[id]) == 0 {
			out = append(out, id)
		}
	}
	return out
}

func (g *Graph) findCycle() string {
	visited := make(map[string]bool, len(g.tasks))
	inStack := make(map[string]bool, len(g.tasks))
	for id := range g.tasks {
		if !visited[id] {
			if node := dfs(id, g.children, visited, inStack); node != "" {
				return node
			}
		}
	}
	return ""
}

func dfs(id string, children map[string][]string, visited, inStack map[string]bool) string {
	visited[id] = true
	inStack[id] = true
	for _, child := range children[id] {
		if !visited[child] {
			if node := dfs(child, children, visited, inStack); node != "" {
				return node
			}
		} else if inStack[child] {
			return child
		}
	}
	inStack[id] = false
	return ""
}
