package ui

import (
    "hash/fnv"
    "github.com/charmbracelet/lipgloss"
)

var tagPalette = []lipgloss.Color{
    "#E06C75", "#E5C07B", "#98C379", "#56B6C2", "#61AFEF",
    "#C678DD", "#D19A66", "#BE5046", "#2BBAC5", "#A9C34F",
}

func TagColor(tag string) lipgloss.Color {
    h := fnv.New32a()
    _, _ = h.Write([]byte(tag))
    return tagPalette[h.Sum32()%uint32(len(tagPalette))]
}

var (
    colorText      = lipgloss.Color("#FFFFFF")
    colorSubtle    = lipgloss.Color("#45475A")
    colorMuted     = lipgloss.Color("#6C7086")
    colorAccent    = lipgloss.Color("#89B4FA")
    colorHighlight = lipgloss.Color("#7D56F4")

    Base = lipgloss.NewStyle().Foreground(colorText)

    Selected = lipgloss.NewStyle().
        Background(colorHighlight).
        Foreground(colorText).
        Bold(true)

    Dimmed = lipgloss.NewStyle().
        Foreground(colorMuted)

    SectionTitle = lipgloss.NewStyle().
        Bold(true).
        Foreground(colorAccent)

    ErrorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#F38BA8")).
        Bold(true)

    DoneStyle = lipgloss.NewStyle().
        Foreground(colorMuted).
        Strikethrough(true)
)

func TagBadge(tag string) string {
    return lipgloss.NewStyle().Foreground(TagColor(tag)).Bold(true).Render(tag)
}

func StatusBadge(status string) string {
    var color lipgloss.Color
    switch status {
    case "todo":       color = "#89B4FA"
    case "in-progress": color = "#A6E3A1"
    case "blocked":    color = "#F38BA8"
    case "done":       color = "#6C7086"
    default:           color = "#FFFFFF"
    }
    return lipgloss.NewStyle().Foreground(color).Render("[" + status + "]")
}

func PriorityBadge(priority string) string {
    var color lipgloss.Color
    switch priority {
    case "high":   color = "#F38BA8"
    case "medium": color = "#E5C07B"
    case "low":    color = "#89B4FA"
    default:       color = "#FFFFFF"
    }
    return lipgloss.NewStyle().Foreground(color).Render("[" + priority + "]")
}
