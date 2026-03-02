package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dayplanner/internal/model"
)

var (
	fieldLabel = lipgloss.NewStyle().
			Width(12).
			Foreground(lipgloss.Color("#89B4FA"))

	fieldActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CDD6F4")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("#89B4FA"))

	fieldInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086"))

	chainItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E3A1"))
)

var formFields = []string{"Tag", "Name", "Priority", "Goal", "DependsOn", "Notes"}

func RenderAddEdit(m model.Model) string {
	var b strings.Builder

	title := "New Task"
	if m.Form.EditingID != "" {
		title = "Edit Task · " + m.Form.EditingID
	}
	b.WriteString(RenderHeader(title))

	if m.Form.IsChain && len(m.Form.ChainIDs) > 0 {
		b.WriteString(renderChainProgress(m.Form.ChainIDs) + "\n")
	}

	for i, name := range formFields {
		value := fieldValue(m.Form, i)
		b.WriteString(renderField(name, value, i == m.Form.FocusIndex))
	}

	if m.Err != nil {
		b.WriteString("\n" + RenderError(m.Err) + "\n")
	}

	b.WriteString(RenderFooter(model.PageAddEdit))
	return b.String()
}

func renderField(label, value string, active bool) string {
	l := fieldLabel.Render(label + ":")
	var v string
	if active {
		v = fieldActive.Render(value + "█")
	} else {
		v = fieldInactive.Render(value)
	}
	return fmt.Sprintf("%s  %s\n", l, v)
}

func renderChainProgress(ids []string) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = chainItem.Render(id)
	}
	return "  Chain: " + strings.Join(parts, Dimmed.Render(" → ")) + Dimmed.Render(" → [next]")
}

func fieldValue(f model.FormState, index int) string {
	if index < len(f.Fields) {
		return f.Fields[index]
	}
	return ""
}
