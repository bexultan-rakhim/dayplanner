package ui

import "dayplanner/internal/model"

func Render(m model.Model) string {
	switch m.Page {
		case model.PageDashboard:
			return RenderDashboard(m)
		case model.PageAddEdit:
			return RenderAddEdit(m)
		case model.PageGraph:
			return RenderGraphView(m)
		case model.PageTaskView:
			return RenderTaskView(m)
		default:
			return RenderDashboard(m)
	}
}
