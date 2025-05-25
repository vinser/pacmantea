package model

import (
	"fmt"
	"strings"

	"github.com/vinser/pacmantea/internal/ui"
	"github.com/vinser/pacmantea/internal/utils"
)

// View function to render entities
func (m *Model) View() string {
	if m.LevelWin {
		if m.GameWin {
			return "You Win! \nPress space to restart. Press 'q' to quit."
		} else {
			view := fmt.Sprintf("Level %d completed! \nPress space to continue. Press 'q' to quit.", m.CurrentLevel+1)
			view += fmt.Sprintf("\nLevel elapsed time: %d seconds!", m.ElapsedTime[m.LevelName])
			view += fmt.Sprintf("\n%v", m.ElapsedTime)
			return view
		}
	}

	if m.GameOver {
		if m.Lives > 1 {
			return fmt.Sprintf("You lost a life! Lives remaining: %d.\nPress space to restart the current level. Press 'q' to quit.", m.Lives-1)
		}
		return "Game Over! \nPress space to restart from the beginning. Press 'q' to quit."
	}
	grid := make([]string, len(m.Maze))
	copy(grid, m.Maze)

	// Place dots
	for _, d := range m.Dots {
		grid[d.Position.Y] = utils.ReplaceAtIndex(grid[d.Position.Y], '·', d.Position.X)
	}

	// Place energizers
	for _, e := range m.Energizers {
		grid[e.Position.Y] = utils.ReplaceAtIndex(grid[e.Position.Y], 'o', e.Position.X)
	}

	// Place ghosts
	for _, g := range m.Ghosts {
		ghostChar := g.Badge
		if g.Dead {
			ghostChar = ' '
		}
		grid[g.Position.Y] = utils.ReplaceAtIndex(grid[g.Position.Y], ghostChar, g.Position.X)
	}

	// Place the pacman with chewing effect
	pacmanChar := 'C'
	if m.Pacman.ChewState {
		pacmanChar = 'c'
	}
	grid[m.Pacman.Position.Y] = utils.ReplaceAtIndex(grid[m.Pacman.Position.Y], pacmanChar, m.Pacman.Position.X)

	// Apply styles to the grid
	for y, row := range grid {
		coloredRow := ""
		for _, rn := range row {
			switch rn {
			case '│', '─', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '║', '═', '╔', '╗', '╚', '╝', '╟', '╢', '╤', '╧', '╖', '╓', '╜', '╙', '╨', '╥':
				coloredRow += ui.WallStyle.Render(string(rn))
			case 'C', 'c':
				coloredRow += renderPacman(m, rn)
			case 'B', 'I', 'P', 'Y':
				coloredRow += renderGhost(m, rn)
			case '.':
				coloredRow += ui.DotStyle.Render(string(rn))
			case 'o':
				coloredRow += ui.EnergyStyle.Render(string(rn))
			default:
				coloredRow += string(rn)
			}
		}
		grid[y] = coloredRow
	}

	// Build the string for display
	view := strings.Join(grid, "\n")
	view += fmt.Sprintf("\nLevel: %d/%d, Score: %d/%d, Lives: %d", m.CurrentLevel+1, len(m.Levels), m.LevelScore, len(m.Dots), m.Lives)
	view += "\nUse arrow keys to move. Press 'q' to quit, 'm' to mute"

	return view
}

func renderPacman(m *Model, r rune) string {
	var rn string
	switch r {
	case 'C':
		rn = m.Config.Badges.Pacman[m.Config.Levels[m.CurrentLevel].PacmanBadge]["open"]
		if m.Pacman.RampantState {
			return ui.EnergyStyle.Render(rn)
		} else {
			return m.Pacman.Style.Render(rn)
		}
	case 'c':
		switch m.Pacman.Move {
		case utils.Direction{X: 1, Y: 0}: // Moving right
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.CurrentLevel].PacmanBadge]["right"]
		case utils.Direction{X: -1, Y: 0}: // Moving left
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.CurrentLevel].PacmanBadge]["left"]
		case utils.Direction{X: 0, Y: -1}: // Moving up
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.CurrentLevel].PacmanBadge]["up"]
		case utils.Direction{X: 0, Y: 1}: // Moving down
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.CurrentLevel].PacmanBadge]["down"]
		default:
			rn = m.Config.Badges.Pacman[m.Config.Levels[m.CurrentLevel].PacmanBadge]["open"]
		}
		if m.Pacman.RampantState {
			if m.Pacman.CooldownState {
				return ui.PacmanStyle.Render(rn)
			} else {
				return ui.EnergyStyle.Render(rn)
			}
		} else {
			return ui.PacmanStyle.Render(rn)
		}
	}
	return ui.PacmanStyle.Render(string(rn))
}

func renderGhost(m *Model, r rune) string {
	rn := m.Config.Badges.Ghosts[m.Config.Levels[m.CurrentLevel].GhostBadges][string(r)]
	switch r {
	case 'B':
		return ui.BlinkyStyle.Render(rn)
	case 'I':
		return ui.InkyStyle.Render(rn)
	case 'P':
		return ui.PinkyStyle.Render(rn)
	case 'Y':
		return ui.ClydeStyle.Render(rn)
	}
	return rn
}
