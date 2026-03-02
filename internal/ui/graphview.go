package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dayplanner/internal/domain"
	"dayplanner/internal/model"
)

const (
	nodeWidth = 20
	connector = " ──► "
)

func RenderGraphView(m model.Model) string {
	var b strings.Builder
	b.WriteString(RenderHeader("Dependency Graph", m.Width))

	if m.Graph == nil || len(m.Tasks) == 0 {
		b.WriteString(Dimmed.Render("  No tasks yet.") + "\n")
		b.WriteString(RenderFooter(model.PageGraph, m.Width))
		return b.String()
	}

	layers := m.Graph.Layers()
	criticalSet := criticalPathSet(m)

	maxRows := 0
	for _, layer := range layers {
		if len(layer) > maxRows {
			maxRows = len(layer)
		}
	}

	for row := 0; row < maxRows; row++ {
		var line strings.Builder
		for col, layer := range layers {
			if row >= len(layer) {
				line.WriteString(strings.Repeat(" ", nodeWidth+len(connector)))
				continue
			}
			id := layer[row]
			task, ok := m.TaskByID(id)
			if !ok {
				continue
			}
			node := renderNode(task, id == m.ActiveTaskID, criticalSet[id])
			if col < len(layers)-1 && hasConnectionRight(layers, row, col, m) {
				line.WriteString(fmt.Sprintf("%-*s%s", nodeWidth, node, connector))
			} else {
				line.WriteString(fmt.Sprintf("%-*s%s", nodeWidth, node, strings.Repeat(" ", len(connector))))
			}
		}
		b.WriteString(line.String() + "\n")
	}

	criticalPath := m.Graph.CriticalPath()
	if len(criticalPath) > 0 {
		b.WriteString("\n" + Dimmed.Render("critical: "+strings.Join(criticalPath, " → ")) + "\n")
	}

	b.WriteString(RenderFooter(model.PageGraph, m.Width))
	return b.String()
}

func renderNode(task domain.Task, isActive, isCritical bool) string {
	label := truncate(task.ID, nodeWidth-2)
	node := "[" + label + "]"

	switch {
	case isActive:
		return Selected.Render(node)
	case isCritical:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F9E2AF")).Bold(true).Render(node)
	case task.Status == domain.StatusDone:
		return Dimmed.Render(node)
	case task.Status == domain.StatusBlocked:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8")).Render(node)
	default:
		return lipgloss.NewStyle().Foreground(TagColor(task.Tag)).Render(node)
	}
}

func criticalPathSet(m model.Model) map[string]bool {
	path := m.Graph.CriticalPath()
	set := make(map[string]bool, len(path))
	for _, id := range path {
		set[id] = true
	}
	return set
}

func hasConnectionRight(layers [][]string, row, col int, m model.Model) bool {
	if col+1 >= len(layers) || row >= len(layers[col]) {
		return false
	}
	leftID := layers[col][row]
	for _, rightID := range layers[col+1] {
		task, ok := m.TaskByID(rightID)
		if !ok {
			continue
		}
		for _, dep := range task.DependsOn {
			if dep == leftID {
				return true
			}
		}
	}
	return false
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
