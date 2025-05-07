package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View function to render entities
func (m model) View() string {
	if m.win {
		if m.winGame {
			return "You Win! \nPress space to restart. Press 'q' to quit."
		} else {
			return fmt.Sprintf("Level %d completed! \nPress space to continue. Press 'q' to quit.", m.currentLevel+1)
		}
	}
	if m.gameOver {
		if m.lives > 1 {
			return fmt.Sprintf("You lost a life! Lives remaining: %d.\nPress space to restart the current level. Press 'q' to quit.", m.lives-1)
		}
		return "Game Over! \nPress space to restart from the beginning. Press 'q' to quit."
	}

	grid := make([]string, len(m.maze))
	copy(grid, m.maze)

	// Place dots
	for _, d := range m.dots {
		grid[d.position.y] = replaceAtIndex(grid[d.position.y], '·', d.position.x)
	}

	// Place energizers
	for _, e := range m.energizers {
		grid[e.position.y] = replaceAtIndex(grid[e.position.y], 'o', e.position.x)
	}

	// Place ghosts
	for _, g := range m.ghosts {
		ghostChar := g.badge
		if g.dead {
			ghostChar = ' '
		}
		grid[g.position.y] = replaceAtIndex(grid[g.position.y], ghostChar, g.position.x)
	}

	// Place the pacman with chewing effect
	pacmanChar := 'C'
	if m.pacman.chewState {
		pacmanChar = 'c'
	}
	grid[m.pacman.position.y] = replaceAtIndex(grid[m.pacman.position.y], pacmanChar, m.pacman.position.x)

	// Apply styles to the grid
	for y, row := range grid {
		coloredRow := ""
		for _, rn := range row {
			switch rn {
			case '│', '─', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '║', '═', '╔', '╗', '╚', '╝', '╟', '╢', '╤', '╧', '╖', '╓', '╜', '╙', '╨', '╥':
				coloredRow += wallStyle.Render(string(rn))
			case 'C', 'c':
				coloredRow += m.renderPacman(rn)
			case 'B', 'I', 'P', 'Y':
				coloredRow += m.renderGhost(rn)
			case '.':
				coloredRow += dotStyle.Render(string(rn))
			case 'o':
				coloredRow += energyStyle.Render(string(rn))
			default:
				coloredRow += string(rn)
			}
		}
		grid[y] = coloredRow
	}

	// Build the string for display
	view := strings.Join(grid, "\n")
	view += fmt.Sprintf("\nLevel: %d/%d, Score: %d/%d, Lives: %d", m.currentLevel+1, len(m.Levels), m.score, m.maxScore, m.lives)
	view += "\nUse arrow keys to move. Press 'q' to quit."

	return view
}

func (m *model) renderPacman(r rune) string {
	var rn string
	switch r {
	case 'C':
		rn = m.Config.Badges.Pacman[m.Config.Levels[m.currentLevel].PacmanBadge]["open"]
		if m.pacman.rampantState {
			return energyStyle.Render(rn)
		} else {
			return m.pacman.style.Render(rn)
		}
	case 'c':
		switch m.pacman.move {
		case direction{x: 1, y: 0}: // Moving right
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.currentLevel].PacmanBadge]["right"]
		case direction{x: -1, y: 0}: // Moving left
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.currentLevel].PacmanBadge]["left"]
		case direction{x: 0, y: -1}: // Moving up
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.currentLevel].PacmanBadge]["up"]
		case direction{x: 0, y: 1}: // Moving down
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.currentLevel].PacmanBadge]["down"]
		default:
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.currentLevel].PacmanBadge]["open"]
		}
		if m.pacman.rampantState {
			if m.pacman.cooldownState {
				return pacmanStyle.Render(rn)
			} else {
				return energyStyle.Render(rn)
			}
		} else {
			return pacmanStyle.Render(rn)
		}
	}
	return pacmanStyle.Render(string(rn))
}

func (m *model) renderGhost(r rune) string {
	var rn string
	switch r {
	case 'B':
		rn = m.Config.Badges.Ghosts[m.Config.Levels[m.currentLevel].GhostBadges][string(r)]
		return blinkyStyle.Render(rn)
	case 'I':
		rn = m.Config.Badges.Ghosts[m.Config.Levels[m.currentLevel].GhostBadges][string(r)]
		return inkyStyle.Render(rn)
	case 'P':
		rn = m.Config.Badges.Ghosts[m.Config.Levels[m.currentLevel].GhostBadges][string(r)]
		return pinkyStyle.Render(rn)
	case 'Y':
		rn = m.Config.Badges.Ghosts[m.Config.Levels[m.currentLevel].GhostBadges][string(r)]
		return clydeStyle.Render(rn)
	}
	return rn
}

// Define styles for different elements
var (
	wallStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))            // Green
	pacmanStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true) // Yellow
	dotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))           // White
	energyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true) // Blue
)

// Define styles for different ghosts
var (
	blinkyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)   // Red
	inkyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)   // Cyan
	pinkyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Bold(true) // Pink
	clydeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Orange
)
