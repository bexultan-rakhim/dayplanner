package history

import "dayplanner/internal/domain"

const maxDepth = 25

type ActionKind int

const (
	ActionDelete ActionKind = iota
	ActionEdit
)

func (a ActionKind) String() string {
	switch a {
	case ActionDelete:
		return "delete"
	case ActionEdit:
		return "edit"
	default:
		return "unknown"
	}
}

type Entry struct {
	Kind     ActionKind
	Snapshot domain.Task
	Index    int
}

type Stack struct {
	undo []Entry
	redo []Entry
}

func New() *Stack {
	return &Stack{
		undo: make([]Entry, 0, maxDepth),
		redo: make([]Entry, 0, maxDepth),
	}
}

func (s *Stack) Push(e Entry) {
	if len(s.undo) >= maxDepth {
		s.undo = s.undo[1:]
	}
	s.undo = append(s.undo, e)
	s.redo = s.redo[:0]
}

func (s *Stack) Undo() (Entry, bool) {
	if len(s.undo) == 0 {
		return Entry{}, false
	}
	last := len(s.undo) - 1
	e := s.undo[last]
	s.undo = s.undo[:last]
	if len(s.redo) >= maxDepth {
		s.redo = s.redo[1:]
	}
	s.redo = append(s.redo, e)
	return e, true
}

func (s *Stack) Redo() (Entry, bool) {
	if len(s.redo) == 0 {
		return Entry{}, false
	}
	last := len(s.redo) - 1
	e := s.redo[last]
	s.redo = s.redo[:last]
	if len(s.undo) >= maxDepth {
		s.undo = s.undo[1:]
	}
	s.undo = append(s.undo, e)
	return e, true
}

func (s *Stack) CanUndo() bool  { return len(s.undo) > 0 }
func (s *Stack) CanRedo() bool  { return len(s.redo) > 0 }
func (s *Stack) UndoDepth() int { return len(s.undo) }
func (s *Stack) RedoDepth() int { return len(s.redo) }

// RestoreIndex returns a safe insertion index given the current task count.
// Use this when handling ActionDelete in the model layer.
func RestoreIndex(e Entry, taskCount int) int {
	if e.Index > taskCount {
		return taskCount
	}
	return e.Index
}
