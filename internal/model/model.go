package model

import (
	"strings"

	"dayplanner/internal/domain"
	"dayplanner/internal/graph"
	"dayplanner/internal/history"
	"dayplanner/internal/repository"
	"dayplanner/internal/scheduler"
)

type Page int

const (
	PageDashboard Page = iota
	PageAddEdit
	PageGraph
	PageTaskView
)

type PickerKind int

const (
	PickerNone     PickerKind = iota
	PickerPriority
	PickerDeps
)

type Picker struct {
	Kind    PickerKind
	Filter  string
	Cursor  int
	Options []string
}

func (p Picker) Filtered() []string {
	if p.Filter == "" {
		return p.Options
	}
	filter := strings.ToLower(p.Filter)
	var out []string
	for _, o := range p.Options {
		if strings.Contains(strings.ToLower(o), filter) {
			out = append(out, o)
		}
	}
	return out
}

type FormField int

const (
	FieldTag      FormField = iota
	FieldName
	FieldPriority
	FieldGoal
	FieldDeps
	FieldNotes
	FieldCount
)

type FormState struct {
	Tag       string
	Name      string
	Priority  domain.Priority
	Goal      string
	DependsOn []string
	Notes     string

	FocusIndex FormField
	Picker     Picker
	ChainIDs   []string
	IsChain    bool
	EditingID  string
}

func (f FormState) PickerOpen() bool {
	return f.Picker.Kind != PickerNone
}

type Model struct {
	Page      Page
	Tasks     []domain.Task
	Graph     *graph.Graph
	Scheduled []domain.Task
	History   *history.Stack
	Repo      repository.Repository

	Cursor       int
	ActiveTaskID string
	Form         FormState

	Width  int
	Height int
	Err    error
}

func New(repo repository.Repository) (Model, error) {
	tasks, err := repo.Load()
	if err != nil {
		return Model{}, err
	}
	m := Model{
		Page:    PageDashboard,
		Tasks:   tasks,
		History: history.New(),
		Repo:    repo,
	}
	return m.RebuildGraph(), nil
}

func (m Model) RebuildGraph() Model {
    g, err := graph.Build(m.Tasks)
    if err != nil {
        m.Err = err
        return m
    }
    m.Graph = g
    m.Tasks = syncBlockedStatus(m.Tasks, g)
    m.Scheduled = scheduler.Order(m.Tasks, g)
    return m
}

func syncBlockedStatus(tasks []domain.Task, g *graph.Graph) []domain.Task {
    doneSet := make(map[string]bool, len(tasks))
    for _, t := range tasks {
        if t.Status == domain.StatusDone {
            doneSet[t.ID] = true
        }
    }

    updated := make([]domain.Task, len(tasks))
    copy(updated, tasks)

    for i, t := range updated {
        if t.Status == domain.StatusInProgress || t.Status == domain.StatusDone {
            continue
        }
        deps := g.Blocking(t.ID)
        isBlocked := false
        for _, dep := range deps {
            if !doneSet[dep] {
                isBlocked = true
                break
            }
        }
        updated[i].SetBlocked(isBlocked)
    }

    return updated
}

func (m Model) TaskByID(id string) (domain.Task, bool) {
	for _, t := range m.Tasks {
		if t.ID == id {
			return t, true
		}
	}
	return domain.Task{}, false
}

func (m Model) SelectedTask() (domain.Task, bool) {
	if m.Cursor < 0 || m.Cursor >= len(m.Scheduled) {
		return domain.Task{}, false
	}
	return m.Scheduled[m.Cursor], true
}

func (m Model) WithPage(p Page) Model {
	m.Page = p
	m.Cursor = 0
	return m
}

func (m Model) WithActiveTask(id string) Model {
	m.ActiveTaskID = id
	return m
}

func (m Model) WithError(err error) Model {
	m.Err = err
	return m
}

func (m Model) ClearError() Model {
	m.Err = nil
	return m
}

func (m Model) OtherTaskIDs(excludeID string) []string {
	out := make([]string, 0, len(m.Tasks))
	for _, t := range m.Tasks {
		if t.ID != excludeID {
			out = append(out, t.ID)
		}
	}
	return out
}
