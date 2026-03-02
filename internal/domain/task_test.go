package domain

import (
	"testing"
)

func TestNewTask_ValidInput(t *testing.T) {
	task, err := NewTask("AUTH-login-flow", "AUTH", "Implement login", PriorityHigh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Status != StatusTodo {
		t.Errorf("expected status todo, got %q", task.Status)
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewTask_InvalidID(t *testing.T) {
	cases := []string{
		"",
		"auth-login",
		"AUTH_login",
		"A-login",
		"TOOLONG-login",
		"AUTH-Login",
		"AUTH-",
	}
	for _, id := range cases {
		_, err := NewTask(id, "AUTH", "name", PriorityLow)
		if err == nil {
			t.Errorf("expected error for id %q", id)
		}
	}
}

func TestNewTask_InvalidTag(t *testing.T) {
	cases := []string{"A", "TOOLONGTAG", "auth", "AU TH"}
	for _, tag := range cases {
		_, err := NewTask("AUTH-slug", tag, "name", PriorityLow)
		if err == nil {
			t.Errorf("expected error for tag %q", tag)
		}
	}
}

func TestAdvanceStatus(t *testing.T) {
	task, _ := NewTask("AUTH-x", "AUTH", "name", PriorityLow)

	task.AdvanceStatus()
	if task.Status != StatusInProgress {
		t.Errorf("expected in-progress, got %q", task.Status)
	}

	task.AdvanceStatus()
	if task.Status != StatusDone {
		t.Errorf("expected done, got %q", task.Status)
	}

	task.AdvanceStatus()
	if task.Status != StatusDone {
		t.Errorf("done should be terminal, got %q", task.Status)
	}
}

func TestAdvanceStatus_FromBlocked(t *testing.T) {
	task, _ := NewTask("AUTH-x", "AUTH", "name", PriorityLow)
	task.Status = StatusBlocked

	task.AdvanceStatus()
	if task.Status != StatusInProgress {
		t.Errorf("expected in-progress from blocked, got %q", task.Status)
	}
}

func TestSetBlocked(t *testing.T) {
	task, _ := NewTask("AUTH-x", "AUTH", "name", PriorityLow)

	task.SetBlocked(true)
	if task.Status != StatusBlocked {
		t.Errorf("expected blocked, got %q", task.Status)
	}

	task.SetBlocked(false)
	if task.Status != StatusTodo {
		t.Errorf("expected todo, got %q", task.Status)
	}
}

func TestSetBlocked_DoesNotAffectInProgress(t *testing.T) {
	task, _ := NewTask("AUTH-x", "AUTH", "name", PriorityLow)
	task.Status = StatusInProgress

	task.SetBlocked(true)
	if task.Status != StatusInProgress {
		t.Errorf("SetBlocked should not override in-progress, got %q", task.Status)
	}
}

func TestIsTerminal(t *testing.T) {
	task, _ := NewTask("AUTH-x", "AUTH", "name", PriorityLow)

	if task.IsTerminal() {
		t.Error("todo task should not be terminal")
	}
	task.Status = StatusDone
	if !task.IsTerminal() {
		t.Error("done task should be terminal")
	}
}
