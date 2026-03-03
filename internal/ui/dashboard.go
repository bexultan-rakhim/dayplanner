package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dayplanner/internal/domain"
	"dayplanner/internal/model"
)

func RenderDashboard(m model.Model) string {
	var b strings.Builder

	b.WriteString(RenderHeader("DayPlanner", m.Width))

	if m.Err != nil {
		b.WriteString(RenderError(m.Err) + "\n")
	}

	ready, blocked, done := partitionTasks(m.Scheduled, m.Graph)

	b.WriteString(renderSection("READY NOW", ready, m, true))
	b.WriteString(renderSection("BLOCKED", blocked, m, true))
	b.WriteString(renderSection("DONE", done, m, true))

	b.WriteString(RenderFooter(model.PageDashboard, m.Width))
	return b.String()
}

func partitionTasks(scheduled []domain.Task, g interface {
	Blocking(id string) []string
}) (ready, blocked, done []domain.Task) {
	for _, t := range scheduled {
		switch t.Status {
		case domain.StatusDone:
			done = append(done, t)
		case domain.StatusBlocked:
			blocked = append(blocked, t)
		default:
			ready = append(ready, t)
		}
	}
	return
}

func renderSection(title string, tasks []domain.Task, m model.Model, showCursor bool) string {
	var b strings.Builder
	b.WriteString(SectionTitle.Render(fmt.Sprintf("── %s (%d)", title, len(tasks))))
	b.WriteString("\n")

	if len(tasks) == 0 {
		b.WriteString(Dimmed.Render("  (empty)") + "\n\n")
		return b.String()
	}

	for _, t := range tasks {
		b.WriteString(renderTaskRow(t, m, showCursor))
	}
	b.WriteString("\n")
	return b.String()
}

func renderTaskRow(t domain.Task, m model.Model, showCursor bool) string {
	selected, _ := m.SelectedTask()
	isSelected := showCursor && selected.ID == t.ID

	const idWidth = 30
	const statusWidth = 15
	const priorityWidth = 10

	slug := idSlug(t.ID)
	if len(slug) > idWidth-8 {
		slug = slug[:(idWidth-8)]
	}
	idStr := TagBadge(t.Tag) + Dimmed.Render("-"+slug)
	idCol := lipgloss.NewStyle().Width(idWidth).MaxWidth(idWidth).Render(idStr)

	statusCol := lipgloss.NewStyle().Width(statusWidth).MaxWidth(statusWidth).Render(StatusBadge(string(t.Status)))
	priorityCol := lipgloss.NewStyle().Width(priorityWidth).MaxWidth(priorityWidth).Render(PriorityBadge(string(t.Priority)))

	name := t.Name
	if t.Status == domain.StatusDone {
		name = DoneStyle.Render(name)
	}

	line := "  " + idCol + statusCol + priorityCol + name

	if isSelected {
		return Selected.Render(line) + "\n"
	}
	return line + "\n"
}

func idSlug(id string) string {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return id
}
