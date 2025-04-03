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
			return fmt.Sprintf("Level %d completed! \nPress space to continue. Press 'q' to quit.", m.currntLevel+1)
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
		grid[d.position.y] = replaceAtIndex(grid[d.position.y], '.', d.position.x)
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

	// Place the player with blinking effect
	playerChar := 'C'
	if m.player.blinkState {
		playerChar = 'c'
	}
	grid[m.player.position.y] = replaceAtIndex(grid[m.player.position.y], playerChar, m.player.position.x)

	// Apply styles to the grid
	for y, row := range grid {
		coloredRow := ""
		for _, char := range row {
			switch char {
			case '│', '─', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '║', '═', '╔', '╗', '╚', '╝', '╟', '╢', '╤', '╧', '╖', '╓', '╜', '╙':
				coloredRow += wallStyle.Render(string(char))
			case 'C':
				if m.player.rampantState {
					coloredRow += energyStyle.Render(string(char))
				} else {
					coloredRow += m.player.style.Render(string(char))
				}
			case 'c':
				if m.player.rampantState {
					if m.player.cooldownState {
						coloredRow += playerStyle.Render(string(char))
					} else {
						coloredRow += energyStyle.Render(string(char))
					}
				} else {
					coloredRow += playerStyle.Render(string(char))
				}
			case 'B':
				coloredRow += blinkyStyle.Render(string(char))
			case 'I':
				coloredRow += inkyStyle.Render(string(char))
			case 'P':
				coloredRow += pinkyStyle.Render(string(char))
			case 'Y':
				coloredRow += clydeStyle.Render(string(char))
			case '.':
				coloredRow += dotStyle.Render(string(char))
			case 'o':
				coloredRow += energyStyle.Render(string(char))
			default:
				coloredRow += string(char)
			}
		}
		grid[y] = coloredRow
	}

	// Build the string for display
	view := strings.Join(grid, "\n")
	view += fmt.Sprintf("\nLevel: %d/%d, Score: %d/%d, Lives: %d", m.currntLevel+1, len(m.Levels), m.score, m.maxScore, m.lives)
	view += "\nUse arrow keys to move. Press 'q' to quit."

	return view
}

// Define styles for different elements
var (
	wallStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))            // Green
	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true) // Yellow
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

// var pacman = []rune{'⭘', '◐', '◓', '◑', '◒'}

// // Ghosts
// var letters = []rune{'B', 'P', 'I', 'Y'}     	// Letters
// var ghosts = []rune{'👿', '👽', '🤖', '👾'}    // Ghosts
// var hebrew = []rune{'ℵ', 'ℶ', 'ℷ', 'ℸ'}      	// Hebrew ghosts
// var greek = []rune{'Α', 'Β', 'Γ', 'Δ'}       	// Greek ghosts
// var control = []rune{'␊', '␋', '␌', '␍'}		// Control ghosts
// var currency = []rune{'$', '€', '£', '¥'}    	// Currency ghosts
// var mathematics = []rune{'∀', '√', '∂', '∫'} 	// Math ghosts
