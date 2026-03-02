package update

import (
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"dayplanner/internal/domain"
	"dayplanner/internal/history"
	"dayplanner/internal/model"
)

func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case SavedMsg:
		m.Err = nil
		return m, nil

	case SaveErrMsg:
		m.Err = msg.Err
		return m, nil

	case tea.KeyMsg:
		km := model.DefaultKeyMap()
		if key.Matches(msg, km.Quit) && !m.Form.PickerOpen() {
			return m, tea.Quit
		}
		if key.Matches(msg, km.Undo) && m.Page != model.PageAddEdit {
			return handleUndo(m)
		}
		if key.Matches(msg, km.Redo) && m.Page != model.PageAddEdit {
			return handleRedo(m)
		}

		switch m.Page {
		case model.PageDashboard:
			return handleDashboard(m, msg, km)
		case model.PageTaskView:
			return handleTaskView(m, msg, km)
		case model.PageAddEdit:
			if m.Form.PickerOpen() {
				return handlePicker(m, msg, km)
			}
			return handleAddEdit(m, msg, km)
		case model.PageGraph:
			return handleGraph(m, msg, km)
		}
	}

	return m, nil
}

func handleDashboard(m model.Model, msg tea.KeyMsg, km model.KeyMap) (model.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, km.Up):
		if m.Cursor > 0 {
			m.Cursor--
		}
	case key.Matches(msg, km.Down):
		if m.Cursor < len(m.Scheduled)-1 {
			m.Cursor++
		}
	case key.Matches(msg, km.Select):
		task, ok := m.SelectedTask()
		if !ok {
			return m, nil
		}
		m.ActiveTaskID = task.ID
		m.Page = model.PageTaskView
	case key.Matches(msg, km.Status):
		return handleAdvanceStatus(m)
	case key.Matches(msg, km.New):
		m = m.WithPage(model.PageAddEdit)
		m.Form = model.FormState{
			Priority: domain.PriorityMedium,
		}
	case key.Matches(msg, km.Delete):
		return handleDelete(m)
	case key.Matches(msg, km.Graph):
		m = m.WithPage(model.PageGraph)
	}
	return m, nil
}

func handleTaskView(m model.Model, msg tea.KeyMsg, km model.KeyMap) (model.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, km.Back):
		m = m.WithPage(model.PageDashboard)
	case key.Matches(msg, km.Status):
		return handleAdvanceStatus(m)
	case key.Matches(msg, km.Edit):
		task, ok := m.TaskByID(m.ActiveTaskID)
		if !ok {
			return m, nil
		}
		m.Page = model.PageAddEdit
		m.Form = taskToForm(task)
	}
	return m, nil
}

func handleAddEdit(m model.Model, msg tea.KeyMsg, km model.KeyMap) (model.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, km.Back):
		m = m.WithPage(model.PageDashboard)
		m.Form = model.FormState{}

	case key.Matches(msg, km.NextField):
		if m.Form.FocusIndex < model.FieldCount-1 {
			m.Form.FocusIndex++
		}

	case key.Matches(msg, km.PrevField):
		if m.Form.FocusIndex > 0 {
			m.Form.FocusIndex--
		}

	case key.Matches(msg, km.Confirm):
		switch m.Form.FocusIndex {
		case model.FieldPriority:
			m = openPriorityPicker(m)
		case model.FieldDeps:
			m = openDepsPicker(m)
		case model.FieldNotes:
			return handleFormSubmit(m)
		default:
			m.Form.FocusIndex++
		}

	case msg.Type == tea.KeyBackspace:
		m = handleBackspace(m)

	case msg.Type == tea.KeyRunes:
		m = appendToActiveField(m, string(msg.Runes))
	}

	return m, nil
}

func handlePicker(m model.Model, msg tea.KeyMsg, km model.KeyMap) (model.Model, tea.Cmd) {
	p := &m.Form.Picker
	filtered := p.Filtered()

	switch {
	case key.Matches(msg, km.Back):
		m.Form.Picker = model.Picker{}

	case key.Matches(msg, km.Up):
		if p.Cursor > 0 {
			p.Cursor--
		}

	case key.Matches(msg, km.Down):
		if p.Cursor < len(filtered)-1 {
			p.Cursor++
		}

	case key.Matches(msg, km.Confirm):
		if len(filtered) == 0 {
			break
		}
		selected := filtered[p.Cursor]
		switch p.Kind {
		case model.PickerPriority:
			m.Form.Priority = domain.Priority(selected)
			m.Form.Picker = model.Picker{}
			m.Form.FocusIndex++
		case model.PickerDeps:
			if !containsStr(m.Form.DependsOn, selected) {
				m.Form.DependsOn = append(m.Form.DependsOn, selected)
			}
			m.Form.Picker = model.Picker{}
		}

	case msg.Type == tea.KeyBackspace:
		if len(p.Filter) > 0 {
			p.Filter = p.Filter[:len(p.Filter)-1]
			p.Cursor = 0
		} else if p.Kind == model.PickerDeps && len(m.Form.DependsOn) > 0 {
			m.Form.DependsOn = m.Form.DependsOn[:len(m.Form.DependsOn)-1]
			m.Form.Picker = model.Picker{}
		}

	case msg.Type == tea.KeyRunes:
		p.Filter += string(msg.Runes)
		p.Cursor = 0
	}

	return m, nil
}

func handleGraph(m model.Model, msg tea.KeyMsg, km model.KeyMap) (model.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, km.Back):
		m = m.WithPage(model.PageDashboard)
	case key.Matches(msg, km.Select):
		if m.ActiveTaskID != "" {
			m.Page = model.PageTaskView
		}
	}
	return m, nil
}

func openPriorityPicker(m model.Model) model.Model {
	m.Form.Picker = model.Picker{
		Kind:    model.PickerPriority,
		Options: []string{"high", "medium", "low"},
	}
	for i, o := range m.Form.Picker.Options {
		if o == string(m.Form.Priority) {
			m.Form.Picker.Cursor = i
			break
		}
	}
	return m
}

func openDepsPicker(m model.Model) model.Model {
	available := m.OtherTaskIDs(m.Form.EditingID)
	filtered := make([]string, 0, len(available))
	for _, id := range available {
		if !containsStr(m.Form.DependsOn, id) {
			filtered = append(filtered, id)
		}
	}
	m.Form.Picker = model.Picker{
		Kind:    model.PickerDeps,
		Options: filtered,
	}
	return m
}

func handleAdvanceStatus(m model.Model) (model.Model, tea.Cmd) {
	var task domain.Task
	var ok bool
	if m.Page == model.PageTaskView {
		task, ok = m.TaskByID(m.ActiveTaskID)
	} else {
		task, ok = m.SelectedTask()
	}
	if !ok {
		return m, nil
	}
	m.History.Push(history.Entry{Kind: history.ActionEdit, Snapshot: task})
	task.AdvanceStatus()
	m = replaceTask(m, task)
	m = m.RebuildGraph()
	return m, SaveCmd(m.Repo, m.Tasks)
}

func handleDelete(m model.Model) (model.Model, tea.Cmd) {
	task, ok := m.SelectedTask()
	if !ok {
		return m, nil
	}
	m.History.Push(history.Entry{
		Kind:     history.ActionDelete,
		Snapshot: task,
		Index:    m.Cursor,
	})
	m.Tasks = removeTask(m.Tasks, task.ID)
	if m.Cursor >= len(m.Tasks) && m.Cursor > 0 {
		m.Cursor--
	}
	m = m.RebuildGraph()
	return m, SaveCmd(m.Repo, m.Tasks)
}

func handleFormSubmit(m model.Model) (model.Model, tea.Cmd) {
	f := m.Form

	if f.EditingID != "" {
		task, ok := m.TaskByID(f.EditingID)
		if !ok {
			return m, nil
		}
		m.History.Push(history.Entry{Kind: history.ActionEdit, Snapshot: task})
		task.Name = f.Name
		task.Tag = f.Tag
		task.Priority = f.Priority
		task.Goal = f.Goal
		task.DependsOn = f.DependsOn
		task.Notes = f.Notes
		task.UpdatedAt = time.Now().UTC()
		m = replaceTask(m, task)
	} else {
		id := f.Tag + "-" + slugify(f.Name)
		task, err := domain.NewTask(id, f.Tag, f.Name, f.Priority)
		if err != nil {
			m.Err = err
			return m, nil
		}
		task.Goal = f.Goal
		task.Notes = f.Notes
		task.DependsOn = f.DependsOn
		if f.IsChain && len(f.ChainIDs) > 0 {
			task.DependsOn = append(task.DependsOn, f.ChainIDs[len(f.ChainIDs)-1])
		}
		m.Tasks = append(m.Tasks, task)

		if f.IsChain {
			m.Form = model.FormState{
				IsChain:  true,
				ChainIDs: append(f.ChainIDs, task.ID),
				Tag:      f.Tag,
				Priority: domain.PriorityMedium,
			}
			m = m.RebuildGraph()
			return m, SaveCmd(m.Repo, m.Tasks)
		}
	}

	m = m.RebuildGraph()
	m.Form = model.FormState{}
	m = m.WithPage(model.PageDashboard)
	return m, SaveCmd(m.Repo, m.Tasks)
}

func handleUndo(m model.Model) (model.Model, tea.Cmd) {
	e, ok := m.History.Undo()
	if !ok {
		return m, nil
	}
	switch e.Kind {
	case history.ActionEdit:
		m = replaceTask(m, e.Snapshot)
	case history.ActionDelete:
		idx := history.RestoreIndex(e, len(m.Tasks))
		m.Tasks = insertTask(m.Tasks, e.Snapshot, idx)
	}
	m = m.RebuildGraph()
	return m, SaveCmd(m.Repo, m.Tasks)
}

func handleRedo(m model.Model) (model.Model, tea.Cmd) {
	e, ok := m.History.Redo()
	if !ok {
		return m, nil
	}
	switch e.Kind {
	case history.ActionEdit:
		task, ok := m.TaskByID(e.Snapshot.ID)
		if !ok {
			return m, nil
		}
		m.History.Push(history.Entry{Kind: history.ActionEdit, Snapshot: task})
		m = replaceTask(m, e.Snapshot)
	case history.ActionDelete:
		m.Tasks = removeTask(m.Tasks, e.Snapshot.ID)
	}
	m = m.RebuildGraph()
	return m, SaveCmd(m.Repo, m.Tasks)
}

func handleBackspace(m model.Model) model.Model {
	switch m.Form.FocusIndex {
	case model.FieldDeps:
		if len(m.Form.DependsOn) > 0 {
			m.Form.DependsOn = m.Form.DependsOn[:len(m.Form.DependsOn)-1]
		}
	default:
		f := activeFieldPtr(&m.Form)
		if f != nil && len(*f) > 0 {
			*f = (*f)[:len(*f)-1]
		}
	}
	return m
}

func appendToActiveField(m model.Model, chars string) model.Model {
	f := activeFieldPtr(&m.Form)
	if f != nil {
		*f += chars
	}
	return m
}

func activeFieldPtr(f *model.FormState) *string {
	switch f.FocusIndex {
	case model.FieldTag:
		return &f.Tag
	case model.FieldName:
		return &f.Name
	case model.FieldGoal:
		return &f.Goal
	case model.FieldNotes:
		return &f.Notes
	}
	return nil
}

func replaceTask(m model.Model, updated domain.Task) model.Model {
	for i, t := range m.Tasks {
		if t.ID == updated.ID {
			m.Tasks[i] = updated
			return m
		}
	}
	return m
}

func removeTask(tasks []domain.Task, id string) []domain.Task {
	out := make([]domain.Task, 0, len(tasks)-1)
	for _, t := range tasks {
		if t.ID != id {
			out = append(out, t)
		}
	}
	return out
}

func insertTask(tasks []domain.Task, task domain.Task, idx int) []domain.Task {
	out := make([]domain.Task, 0, len(tasks)+1)
	out = append(out, tasks[:idx]...)
	out = append(out, task)
	out = append(out, tasks[idx:]...)
	return out
}

func taskToForm(task domain.Task) model.FormState {
	return model.FormState{
		EditingID: task.ID,
		Tag:       task.Tag,
		Name:      task.Name,
		Priority:  task.Priority,
		Goal:      task.Goal,
		DependsOn: append([]string{}, task.DependsOn...),
		Notes:     task.Notes,
	}
}

func containsStr(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

func slugify(name string) string {
	out := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z':
			out = append(out, c)
		case c >= 'A' && c <= 'Z':
			out = append(out, c+32)
		case c == ' ' || c == '-':
			if len(out) > 0 && out[len(out)-1] != '-' {
				out = append(out, '-')
			}
		}
	}
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}
	return string(out)
}
