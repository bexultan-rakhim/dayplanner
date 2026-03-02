package ui

import (
	"strings"

	"dayplanner/internal/model"
)

func RenderHeader(title string) string {
	return Header.Render(title)
}

func RenderFooter(page model.Page) string {
	keys := footerKeys(page)
	return Footer.Render(strings.Join(keys, "  "))
}

func RenderError(err error) string {
	if err == nil {
		return ""
	}
	return ErrorStyle.Render("error: " + err.Error())
}

func footerKeys(page model.Page) []string {
	switch page {
	case model.PageDashboard:
		return []string{
			"↑/↓ navigate",
			"enter open",
			"space status",
			"n new",
			"g graph",
			"G group",
			"ctrl+z undo",
			"q quit",
		}
	case model.PageAddEdit:
		return []string{
			"tab next field",
			"shift+tab prev",
			"ctrl+r regenerate",
			"enter confirm",
			"esc back",
		}
	case model.PageGraph:
		return []string{
			"↑/↓/←/→ navigate",
			"enter open task",
			"esc back",
		}
	case model.PageTaskView:
		return []string{
			"e edit",
			"space status",
			"enter open dep",
			"esc back",
		}
	default:
		return nil
	}
}
