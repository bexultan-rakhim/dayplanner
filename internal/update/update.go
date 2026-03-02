package update

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"dayplanner/internal/model"
	// "dayplanner/internal/history"
	// "dayplanner/internal/domain"
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

		if key.Matches(msg, km.Quit){
			return m, tea.Quit
		}
		// do nothing for now.
	}

	return m, nil
}
