package ui

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
    "dayplanner/internal/model"
)

func RenderHeader(title string, width int) string {
    titleStr := lipgloss.NewStyle().
        Bold(true).
        Foreground(colorText).
        Render(title)
    sep := lipgloss.NewStyle().
        Foreground(colorSubtle).
        Render(strings.Repeat("─", width))
    return titleStr + "\n" + sep + "\n\n"
}

func RenderFooter(page model.Page, width int) string {
    hints := footerHints(page)
    sep := lipgloss.NewStyle().
        Foreground(colorSubtle).
        Render(strings.Repeat("─", width))

    keyStyle   := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
    labelStyle := lipgloss.NewStyle().Foreground(colorMuted)
    dotStyle   := lipgloss.NewStyle().Foreground(colorSubtle)

    parts := make([]string, len(hints))
    for i, h := range hints {
        parts[i] = keyStyle.Render(h[0]) + " " + labelStyle.Render(h[1])
    }
    line := strings.Join(parts, "  " + dotStyle.Render("·") + "  ")

    return "\n" + sep + "\n" + lipgloss.NewStyle().
        Foreground(colorMuted).
        Render(line)
}

func RenderError(err error) string {
    if err == nil {
        return ""
    }
    return ErrorStyle.Render("✗ " + err.Error())
}

// footerHints returns [key, label] pairs per page.
func footerHints(page model.Page) [][2]string {
    switch page {
    case model.PageDashboard:
        return [][2]string{
            {"j/k", "navigate"},
            {"enter", "open"},
            {"space", "status"},
            {"n", "new"},
            {"g", "graph"},
            {"u", "undo"},
            {"q", "quit"},
        }
    case model.PageAddEdit:
        return [][2]string{
            {"tab", "next field"},
            {"shift+tab", "prev"},
            {"enter", "confirm"},
            {"esc", "back"},
        }
    case model.PageGraph:
        return [][2]string{
            {"j/k", "navigate"},
            {"enter", "open"},
            {"esc", "back"},
        }
    case model.PageTaskView:
        return [][2]string{

            {"e", "edit"},
            {"space", "status"},
            {"esc", "back"},
        }
    default:
        return nil
    }
}
