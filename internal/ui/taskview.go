package ui

import (
	"fmt"
	"strings"

	"dayplanner/internal/model"
)

func RenderTaskView(m model.Model) string {
	var b strings.Builder

	b.WriteString(RenderHeader("Task", m.Width))

	task, ok := m.TaskByID(m.ActiveTaskID)
	if !ok {
		b.WriteString(ErrorStyle.Render("task not found: " + m.ActiveTaskID))
		b.WriteString(RenderFooter(model.PageTaskView, m.Width))
		return b.String()
	}

	titleLine := fmt.Sprintf("%s  %s  %s",
		TagBadge(task.Tag)+Dimmed.Render("-"+idSlug(task.ID)),
		StatusBadge(string(task.Status)),
		PriorityBadge(string(task.Priority)),
	)
	b.WriteString(titleLine + "\n")
	b.WriteString(strings.Repeat("─", 60) + "\n\n")

	b.WriteString(field("Name", task.Name))
	b.WriteString(field("Goal", task.Goal))
	if task.Notes != "" {
		b.WriteString(field("Notes", task.Notes))
	}
	b.WriteString(field("Created", task.CreatedAt.Format("2006-01-02 15:04")))
	b.WriteString(field("Updated", task.UpdatedAt.Format("2006-01-02 15:04")))

	b.WriteString("\n")

	blockedBy := m.Graph.Blocking(task.ID)
	blocks := m.Graph.BlockedBy(task.ID)

	b.WriteString(SectionTitle.Render("── Dependencies") + "\n")

	if len(blockedBy) == 0 {
		b.WriteString(Dimmed.Render("  Blocked by:  (none)") + "\n")
	} else {
		for _, depID := range blockedBy {
			dep, ok := m.TaskByID(depID)
			if !ok {
				continue
			}
			b.WriteString(fmt.Sprintf("  Blocked by:  %s  %s\n",
				TagBadge(dep.Tag)+Dimmed.Render("-"+idSlug(dep.ID)),
				StatusBadge(string(dep.Status)),
			))
		}
	}

	b.WriteString("\n")

	if len(blocks) == 0 {
		b.WriteString(Dimmed.Render("  Blocks:      (none)") + "\n")
	} else {
		for _, depID := range blocks {
			dep, ok := m.TaskByID(depID)
			if !ok {
				continue
			}
			b.WriteString(fmt.Sprintf("  Blocks:      %s  %s\n",
				TagBadge(dep.Tag)+Dimmed.Render("-"+idSlug(dep.ID)),
				StatusBadge(string(dep.Status)),
			))
		}
	}

	b.WriteString(RenderFooter(model.PageTaskView, m.Width))
	return b.String()
}

func field(label, value string) string {
	return fmt.Sprintf("  %-10s %s\n", label+":", value)
}
