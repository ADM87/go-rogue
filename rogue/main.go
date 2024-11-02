package main

import (
	"fmt"
	"os"
	game "rogue/game"

	tea "github.com/charmbracelet/bubbletea"
)

func getFlag(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}

func main() {
	if getFlag("-game") {
		runApplication(game.NewModel())
	}
	fmt.Fprintln(os.Stderr, "Must specify application type")
	os.Exit(1)
}

func runApplication(app tea.Model) {
	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting application: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
