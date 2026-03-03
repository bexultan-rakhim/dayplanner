package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)


type Status string
type Priority string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusBlocked    Status = "blocked"
	StatusDone       Status = "done"
)

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

var validTag = regexp.MustCompile(`^[A-Z]{2,6}$`)
var validID = regexp.MustCompile(`^[A-Z]{2,6}-[a-z0-9]+(-[a-z0-9]+)*$`)

type Task struct {
	ID        string    `json:"id"`
	Tag       string    `json:"tag"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Priority  Priority  `json:"priority"`
	Status    Status    `json:"status"`
	DependsOn []string  `json:"depends_on"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTask(id, tag, name string, priority Priority) (Task, error) {
	if err := validateID(id); err != nil {
		return Task{}, err
	}

	if err := validateTag(tag); err != nil {
		return Task{}, err
	}

	if strings.TrimSpace(name) == ""  {
		return Task{}, fmt.Errorf("name cannot be empty")
	}

	if err := validatePriority(priority); err != nil {
		return Task{}, err
	}

	now := time.Now().UTC()
	return Task {
		ID: id,
		Tag: tag,
		Name: name,
		Priority: priority,
		Status: StatusTodo,
		DependsOn: []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func validateID(id string) error {
	if !validID.MatchString(id) {
		return fmt.Errorf("invalid id %q: must match TAG-slug (e.g. AUTH-login-flow)", id)
	}

	return nil
}

func validateTag(tag string) error {
  if !validTag.MatchString(tag) {
		return fmt.Errorf("tag %q must be 2–6 uppercase letters", tag)
	}

	return nil
}

func validatePriority(p Priority) error {
	switch p {
	case PriorityHigh, PriorityMedium, PriorityLow:
		return nil
	}

	return fmt.Errorf("invalid priority %q", p)
}

func (t Task) IsTerminal() bool {
	return t.Status == StatusDone
}

func (t *Task) AdvanceStatus() {
	switch t.Status {
	case StatusTodo:
		t.Status = StatusInProgress
	case StatusInProgress:
		t.Status = StatusDone
	case StatusDone:
		t.Status = StatusTodo
	case StatusBlocked:
		// blocked status is managed automatically; ignore manual toggle
		return
	}
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) SetBlocked (blocked bool) {
	if blocked  && t.Status == StatusTodo {
		t.Status = StatusBlocked
	} else if !blocked && t.Status == StatusBlocked {
		t.Status = StatusTodo
	}
	t.UpdatedAt = time.Now().UTC()
}
