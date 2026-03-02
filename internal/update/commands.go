package update

import (
	tea "github.com/charmbracelet/bubbletea"
	"dayplanner/internal/repository"
	"dayplanner/internal/domain"
)

type SavedMsg struct{}

type SaveErrMsg struct {
	Err error
}

func SaveCmd(repo repository.Repository, tasks []domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := repo.Save(tasks); err != nil {
			return SaveErrMsg{Err: err}
		}
		return SavedMsg{}
	}
}

