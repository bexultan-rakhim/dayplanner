package model

import (
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

type FormState struct {
	Fields     []string
	FocusIndex int
	ChainIDs   []string
	IsChain    bool
	EditingID  string
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

	Err error
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
	m.Scheduled = scheduler.Order(m.Tasks, g)
	return m
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
