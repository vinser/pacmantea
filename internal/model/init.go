package model

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.pacmanBlinkTick(), // Start the timer for Pac-Man blinking
		m.ghostMoveTick(),   // Start the timer for ghost movement
		splashScreen(),
	}
	return tea.Batch(cmds...)
}
