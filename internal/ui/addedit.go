package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dayplanner/internal/domain"
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

	pickerOverlay = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#89B4FA")).
			Padding(0, 1)

	pickerSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("#313244")).
			Foreground(lipgloss.Color("#CDD6F4")).
			Bold(true)

	pickerNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086"))

	chainItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E3A1"))

	depTag = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6E3A1"))

	depRemove = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F38BA8"))

	fieldActiveBorder = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(lipgloss.Color("#89B4FA"))
)

var priorities = []string{
	string(domain.PriorityHigh),
	string(domain.PriorityMedium),
	string(domain.PriorityLow),
}

func RenderAddEdit(m model.Model) string {
	var b strings.Builder

	title := "New Task"
	if m.Form.EditingID != "" {
		title = "Edit · " + m.Form.EditingID
	}
	if m.Form.IsChain {
		title += "  (chain mode)"
	}
	b.WriteString(RenderHeader(title, m.Width))

	if m.Form.IsChain && len(m.Form.ChainIDs) > 0 {
		b.WriteString(renderChainProgress(m.Form.ChainIDs) + "\n\n")
	}

	editing := m.Form.FieldEditing
	b.WriteString(renderFormField("Tag", m.Form.Tag, m.Form.FocusIndex == model.FieldTag, false, editing))
	b.WriteString(renderFormField("Name", m.Form.Name, m.Form.FocusIndex == model.FieldName, false, editing))
	b.WriteString(renderPriorityField(m))
	b.WriteString(renderFormField("Goal", m.Form.Goal, m.Form.FocusIndex == model.FieldGoal, false, editing))
	b.WriteString(renderDepsField(m))
	b.WriteString(renderFormField("Notes", m.Form.Notes, m.Form.FocusIndex == model.FieldNotes, false, editing))

	if editing {
		b.WriteString("\n" + Dimmed.Render("  esc · stop editing") + "\n")
	} else {
		b.WriteString("\n" + Dimmed.Render("  j/k · navigate   i · edit/open   enter · save") + "\n")
	}

	if m.Form.PickerOpen() {
		b.WriteString("\n" + renderPicker(m.Form.Picker))
	}

	if m.Err != nil {
		b.WriteString("\n" + RenderError(m.Err) + "\n")
	}

	b.WriteString(RenderFooter(model.PageAddEdit, m.Width))
	return b.String()
}

func renderFormField(label, value string, active, readonly, editing bool) string {
	l := fieldLabel.Render(label + ":")
	var v string
	switch {
	case readonly:
		v = fieldInactive.Render(value)
	case active && editing:
		v = fieldActive.Render(value + "█")
	case active:
		v = fieldActive.Render(value)
	default:
		v = fieldInactive.Render(value)
	}
	return fmt.Sprintf("%s  %s\n", l, v)
}

func renderPriorityField(m model.Model) string {
	active := m.Form.FocusIndex == model.FieldPriority
	value := string(m.Form.Priority)
	if value == "" {
		value = "(select)"
	}
	badge := PriorityBadge(value)
	l := fieldLabel.Render("Priority:")
	indicator := ""
	if active {
		badge = fieldActiveBorder.Render(badge)
		indicator = Dimmed.Render("  i · open")
	}
	return fmt.Sprintf("%s  %s%s\n", l, badge, indicator)
}

func renderDepsField(m model.Model) string {
	active := m.Form.FocusIndex == model.FieldDeps
	l := fieldLabel.Render("DependsOn:")

	var depParts []string
	for _, id := range m.Form.DependsOn {
		depParts = append(depParts, depTag.Render(id)+depRemove.Render(" ✕"))
	}

	var value string
	if len(depParts) == 0 {
		value = fieldInactive.Render("(none)")
	} else {
		value = strings.Join(depParts, "  ")
	}

	if active {
		value = fieldActiveBorder.Render(value)
	}

	indicator := ""
	if active {
		indicator = Dimmed.Render("  i · add · backspace to remove last")
	}

	return fmt.Sprintf("%s  %s%s\n", l, value, indicator)
}

func renderPicker(p model.Picker) string {
	filtered := p.Filtered()
	if len(filtered) == 0 {
		return pickerOverlay.Render(Dimmed.Render("  no matches"))
	}

	var rows []string
	if p.Filter != "" {
		rows = append(rows, Dimmed.Render("  filter: ")+p.Filter)
	}

	for i, opt := range filtered {
		var line string
		if i == p.Cursor {
			line = pickerSelected.Render(fmt.Sprintf("▶ %s", opt))
		} else {
			line = pickerNormal.Render(fmt.Sprintf("  %s", opt))
		}
		rows = append(rows, line)
	}

	return pickerOverlay.Render(strings.Join(rows, "\n"))
}

func renderChainProgress(ids []string) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = chainItem.Render(id)
	}
	return "  Chain: " + strings.Join(parts, Dimmed.Render(" → ")) + Dimmed.Render(" → [next]")
}
