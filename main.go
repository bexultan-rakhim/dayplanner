package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"dayplanner/config"
	"dayplanner/internal/model"
	"dayplanner/internal/repository"
)

func main() {
	dir := config.Dir()
	repo, err := repository.NewJSONRepository(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dayplanner: failed to initialise storage: %v\n", err)
		os.Exit(1)
	}

	m, err := model.New(repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dayplanner: failed to load tasks: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		app{m: m},
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "dayplanner: %v\n", err)
		os.Exit(1)
	}
}

