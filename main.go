package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// ghostsLove = true

	model := newModel()
	model.playSound(SOUND_BEGINNING)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
