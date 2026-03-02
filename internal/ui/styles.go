package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

var tagPalette = []lipgloss.Color{
	"#E06C75", // red
	"#E5C07B", // yellow
	"#98C379", // green
	"#56B6C2", // cyan
	"#61AFEF", // blue
	"#C678DD", // purple
	"#D19A66", // orange
	"#BE5046", // dark red
	"#2BBAC5", // teal
	"#A9C34F", // lime
}

func TagColor(tag string) lipgloss.Color {
	h := fnv.New32a()
	_, _ = h.Write([]byte(tag))
	return tagPalette[h.Sum32()%uint32(len(tagPalette))]
}

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	muted     = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}

	Base = lipgloss.NewStyle()

	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#CDD6F4")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(subtle.Dark)).
		MarginBottom(1)

	Footer = lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted.Dark)).
		MarginTop(1)

	Selected = lipgloss.NewStyle().
			Background(lipgloss.Color(highlight.Dark)).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	Dimmed = lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted.Dark))

	SectionTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89B4FA")).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F38BA8")).
			Bold(true)

	DoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(muted.Dark)).
			Strikethrough(true)
)

func TagBadge(tag string) string {
	return lipgloss.NewStyle().
		Foreground(TagColor(tag)).
		Bold(true).
		Render(tag)
}

func StatusBadge(status string) string {
	var color lipgloss.Color
	switch status {
	case "todo":
		color = "#89B4FA"
	case "in-progress":
		color = "#A6E3A1"
	case "blocked":
		color = "#F38BA8"
	case "done":
		color = "#6C7086"
	default:
		color = "#CDD6F4"
	}
	return lipgloss.NewStyle().
		Foreground(color).
		Render("[" + status + "]")
}

func PriorityBadge(priority string) string {
	var color lipgloss.Color
	switch priority {
	case "high":
		color = "#F38BA8"
	case "medium":
		color = "#E5C07B"
	case "low":
		color = "#89B4FA"
	default:
		color = "#CDD6F4"
	}
	return lipgloss.NewStyle().
		Foreground(color).
		Render("[" + priority + "]")
}
