package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type point struct {
	x, y int
}

type direction struct {
	x, y int
}

// Base entity structure
type entity struct {
	position point
	move     direction
	style    lipgloss.Style
	name     string
	badge    rune
}

// Player structure
type player struct {
	entity
	blinkState    bool // For blinking effect
	rampantState  bool // For rampant mode
	cooldownState bool // For cooldown
}

// Ghost structure
type ghost struct {
	entity
}

// Dot structure
type dot struct {
	entity
}

// Energizer structure
type energizer struct {
	entity
}

// Update the model to use the new structures
type model struct {
	maze         []string    // Maze as an array of strings (current state)
	originalMaze []string    // Original maze (immutable)
	maxScore     int         // Original number of dots
	player       player      // Pac-Man
	dots         []dot       // List of dots
	energizers   []energizer // List of energizers
	ghosts       []ghost     // List of ghosts
	score        int         // Score
	gameOver     bool        // Game over flag
	win          bool        // Win flag
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

// Update the initialModel function
func initialModel(maze []string) model {
	rand.Seed(time.Now().UnixNano())

	// Ensure the maze has a minimum size of 5x5
	if len(maze) < 5 || len(maze[0]) < 5 {
		panic("The maze must be at least 5x5")
	}

	var playerEntity player
	dots := []dot{}
	energizers := []energizer{}
	ghosts := []ghost{}

	// Create a copy of the original maze
	originalMaze := make([]string, len(maze))
	copy(originalMaze, maze)

	for y, row := range maze {
		for x, char := range row {
			switch char {
			case 'C':
				playerEntity = player{
					entity: entity{
						position: point{x: x, y: y},
						style:    playerStyle,
						name:     "Pac-Man",
						badge:    'C',
					},
					blinkState: false,
				}
				dots = append(dots, dot{
					entity: entity{
						position: point{x: x, y: y},
						style:    dotStyle,
						name:     "Dot",
						badge:    '.',
					},
				})
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case '.':
				dots = append(dots, dot{
					entity: entity{
						position: point{x: x, y: y},
						style:    dotStyle,
						name:     "Dot",
						badge:    '.',
					},
				})
			case 'o':
				energizers = append(energizers, energizer{
					entity: entity{
						position: point{x: x, y: y},
						style:    energyStyle,
						name:     "Energizer",
						badge:    'o',
					},
				})
			case 'B':
				ghosts = append(ghosts, ghost{
					entity: entity{
						position: point{x: x, y: y},
						style:    blinkyStyle,
						name:     "Blinky",
						badge:    'B',
					},
				})
				dots = append(dots, dot{
					entity: entity{
						position: point{x: x, y: y},
						style:    dotStyle,
						name:     "Dot",
						badge:    '.',
					},
				})
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'I':
				ghosts = append(ghosts, ghost{
					entity: entity{
						position: point{x: x, y: y},
						style:    inkyStyle,
						name:     "Inky",
						badge:    'I',
					},
				})
				dots = append(dots, dot{
					entity: entity{
						position: point{x: x, y: y},
						style:    dotStyle,
						name:     "Dot",
						badge:    '.',
					},
				})
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'P':
				ghosts = append(ghosts, ghost{
					entity: entity{
						position: point{x: x, y: y},
						style:    pinkyStyle,
						name:     "Pinky",
						badge:    'P',
					},
				})
				dots = append(dots, dot{
					entity: entity{
						position: point{x: x, y: y},
						style:    dotStyle,
						name:     "Dot",
						badge:    '.',
					},
				})
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'Y':
				ghosts = append(ghosts, ghost{
					entity: entity{
						position: point{x: x, y: y},
						style:    clydeStyle,
						name:     "Clyde",
						badge:    'Y',
					},
				})
				dots = append(dots, dot{
					entity: entity{
						position: point{x: x, y: y},
						style:    dotStyle,
						name:     "Dot",
						badge:    '.',
					},
				})
				maze[y] = replaceAtIndex(maze[y], '.', x)
			}
		}
	}

	// Convert maze walls to pseudographics for the current maze only
	maze = replaceWallsWithGraphics(maze)

	return model{
		maze:         maze,
		originalMaze: originalMaze, // Keep originalMaze unchanged
		maxScore:     len(dots),
		player:       playerEntity,
		dots:         dots,
		energizers:   energizers,
		ghosts:       ghosts,
		score:        0,
		gameOver:     false,
		win:          false,
	}
}

func (m model) Init() tea.Cmd {
	// Start the timer for ghost movement and blinking
	return tea.Batch(ghostMoveTick(), playerBlinkTick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver || m.win {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				// Restart the game using originalMaze
				newModel := initialModel(m.originalMaze)
				// Start the timer for ghost movement and blinking
				return newModel, tea.Batch(ghostMoveTick(), playerBlinkTick())
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.canMove(m.player.position.x, m.player.position.y-1) {
				m.player.position.y--
				m.player.move = direction{x: 0, y: -1}
			}
		case "down":
			if m.canMove(m.player.position.x, m.player.position.y+1) {
				m.player.position.y++
				m.player.move = direction{x: 0, y: 1}
			}
		case "left":
			newX := m.tunnelMove(m.player.position.x - 1)
			if m.canMove(newX, m.player.position.y) {
				m.player.position.x = newX
				m.player.move = direction{x: -1, y: 0}
			}
		case "right":
			newX := m.tunnelMove(m.player.position.x + 1)
			if m.canMove(newX, m.player.position.y) {
				m.player.position.x = newX
				m.player.move = direction{x: 1, y: 0}
			}
		}

		// Check for dot collection
		for i := len(m.dots) - 1; i >= 0; i-- {
			if m.player.position == m.dots[i].position {
				m.score++
				m.maze[m.player.position.y] = replaceAtIndex(m.maze[m.player.position.y], ' ', m.player.position.x) // Replace dot with a space
				m.dots = append(m.dots[:i], m.dots[i+1:]...)
			}
		}

		// Check for win condition
		if len(m.dots) == 0 {
			m.win = true
		}

		// Check for energizer collection
		for i := len(m.energizers) - 1; i >= 0; i-- {
			if m.player.position == m.energizers[i].position {
				m.score++
				m.maze[m.player.position.y] = replaceAtIndex(m.maze[m.player.position.y], ' ', m.player.position.x) // Replace dot with a space
				m.energizers = append(m.energizers[:i], m.energizers[i+1:]...)

				// Activate rampant mode
				m.player.rampantState = true
				return m, startRampantTimer()
			}
		}

		// Check for collision with ghosts
		for _, ghost := range m.ghosts {
			if m.player.position == ghost.position {
				m.gameOver = true
			}
		}

	case ghostMoveMsg:
		// Move ghosts
		for i := range m.ghosts {
			if m.player.rampantState {
				m.ghosts[i].position = m.escapeMove(m.ghosts[i].position.x, m.ghosts[i].position.y)
				continue
			}
			switch m.ghosts[i].name {
			case "Blinky":
				m.ghosts[i].position = m.straitMove(m.ghosts[i].position.x, m.ghosts[i].position.y)
			case "Inky":
				m.ghosts[i].position = m.chaosMove(m.ghosts[i].position.x, m.ghosts[i].position.y)
			case "Pinky":
				m.ghosts[i].position = m.predictMove(m.ghosts[i].position.x, m.ghosts[i].position.y)
			case "Clyde":
				m.ghosts[i].position = m.cagyMove(m.ghosts[i].position.x, m.ghosts[i].position.y)
			}
		}

		// Check for collision with ghosts after their movement
		for _, ghost := range m.ghosts {
			if m.player.position == ghost.position {
				m.gameOver = true
			}
		}

		// Start the next tick for ghost movement
		return m, ghostMoveTick()

	case playerBlinkMsg:
		// Toggle the blink state
		m.player.blinkState = !m.player.blinkState
		// Schedule the next blink
		return m, playerBlinkTick()

	case rampantEndMsg:
		// Start cooldown
		m.player.cooldownState = true
		return m, startCooldownTimer()

	case cooldownEndMsg:
		// End cooldown and fully reset Pac-Man's state
		m.player.rampantState = false
		m.player.cooldownState = false
		return m, nil
	}

	return m, nil
}

func (m model) chaosMove(x, y int) point {
	directions := randomDirections()
	return m.ghostMove(x, y, directions)
}

func (m model) straitMove(x, y int) point {
	directions := sortDirectionsByDistance(point{x, y}, m.player.position)
	return m.ghostMove(x, y, directions)
}

func (m model) predictMove(x, y int) point {
	destination := point{x: m.player.position.x + m.player.move.x, y: m.player.position.y + m.player.move.y}
	directions := sortDirectionsByDistance(point{x, y}, destination)
	return m.ghostMove(x, y, directions)
}

func (m model) cagyMove(x, y int) point {
	destination := point{x: m.player.position.x - 2*m.player.move.x, y: m.player.position.y - 2*m.player.move.y}
	directions := sortDirectionsByDistance(point{x, y}, destination)
	return m.ghostMove(x, y, directions)
}

func (m model) escapeMove(x, y int) point {
	destination := point{x: 2*x - m.player.position.x, y: 2*y - m.player.position.y}
	directions := sortDirectionsByDistance(point{x, y}, destination)
	return m.ghostMove(x, y, directions)
}

func (m model) ghostMove(x, y int, directions []point) point {
	for _, dir := range directions {
		newX, newY := x+dir.x, y+dir.y
		newX = m.tunnelMove(newX)
		if m.canMove(newX, newY) {
			return point{x: newX, y: newY}
		}
	}
	return point{x: x, y: y}
}

// Sort directions by distance to the player
func sortDirectionsByDistance(p1, p2 point) []point {
	directions := []point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	sort.Slice(directions, func(i, j int) bool {
		d1 := distance(point{x: p1.x + directions[i].x, y: p1.y + directions[i].y}, p2)
		d2 := distance(point{x: p1.x + directions[j].x, y: p1.y + directions[j].y}, p2)
		return d1 < d2
	})
	return directions
}

func distance(p1, p2 point) float64 {
	return math.Sqrt(float64((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y)))
}

func randomDirections() []point {
	directions := []point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	return directions
}

func (m model) tunnelMove(newX int) int {
	if newX < 0 {
		return len([]rune(m.maze[0])) - 1
	}
	if newX > len([]rune(m.maze[0]))-1 {
		return 0
	}
	return newX
}

// Update the View function to render entities
func (m model) View() string {
	if m.win {
		return "You Win! Press space to restart.\nPress 'q' to quit."
	}
	if m.gameOver {
		return "Game Over! Press space to restart.\nPress 'q' to quit."
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
		grid[g.position.y] = replaceAtIndex(grid[g.position.y], g.badge, g.position.x)
	}

	// Place the player with blinking effect
	playerChar := 'C'
	if m.player.blinkState {
		playerChar = 'c'
	}
	grid[m.player.position.y] = replaceAtIndex(grid[m.player.position.y], playerChar, m.player.position.x)

	// Apply colors to the grid
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
	view += fmt.Sprintf("\nScore: %d/%d\n", m.score, m.maxScore)
	view += "Use arrow keys to move. Press 'q' to quit."

	return view
}

// Helper function to replace a character in a string at a given index
func replaceAtIndex(s string, r rune, index int) string {
	runes := []rune(s)
	runes[index] = r
	return string(runes)
}

// Check if movement is possible
func (m model) canMove(x, y int) bool {

	// Convert the maze row to runes to handle multi-byte characters
	row := []rune(m.maze[y])
	char := row[x]

	// List of all wall characters (pseudographics)
	wallChars := []rune{'│', '─', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '║', '═', '╔', '╗', '╚', '╜', '╟', '╢', '╤', '╧', '╖', '╓', '╜', '╙'}

	// Check if the character is a wall
	for _, wallChar := range wallChars {
		if char == wallChar {
			return false
		}
	}

	return true // Not a wall, movement is allowed
}

// Message type for ghost movement
type ghostMoveMsg struct{}

const ghostMoveTickDuration = time.Second / 2

// Command to trigger ghost movement ticks
func ghostMoveTick() tea.Cmd {
	return tea.Tick(ghostMoveTickDuration, func(_ time.Time) tea.Msg {
		return ghostMoveMsg{}
	})
}

// Message type for blinking
type playerBlinkMsg struct{}

const playerBlinkTickDuration = time.Second / 2

// Command to trigger Pac-Man blinking
func playerBlinkTick() tea.Cmd {
	return tea.Tick(playerBlinkTickDuration, func(_ time.Time) tea.Msg {
		return playerBlinkMsg{}
	})
}

// Message types for rampant and cooldown states
type rampantEndMsg struct{}

const rampantTimerDuration = 5 * time.Second

// Command to start the rampant timer
func startRampantTimer() tea.Cmd {
	return tea.Tick(rampantTimerDuration, func(_ time.Time) tea.Msg {
		return rampantEndMsg{}
	})
}

type cooldownEndMsg struct{}

const cooldownTimerDuration = 2 * time.Second

// Command to start the cooldown timer
func startCooldownTimer() tea.Cmd {
	return tea.Tick(cooldownTimerDuration, func(_ time.Time) tea.Msg {
		return cooldownEndMsg{}
	})
}

// Convert maze walls to pseudographics
func replaceWallsWithGraphics(maze []string) []string {
	height := len(maze)
	width := len(maze[0])

	// Fix the outer walls tunnels
	for i, line := range maze {
		row := []rune(line)
		switch {
		case row[0] == ' ' && row[width-1] != ' ':
			row[width-1] = ' '
			maze[i] = string(row)
		case row[0] != ' ' && row[width-1] == ' ':
			row[0] = ' '
			maze[i] = string(row)
		}
	}

	// Create a new grid for pseudographics
	newMaze := make([]string, height)
	for y := 0; y < height; y++ {
		newRow := []rune(maze[y])
		for x := 0; x < width; x++ {
			if maze[y][x] == '#' {
				// Determine neighbors of the current cell
				top := y > 0 && maze[y-1][x] == '#'
				bottom := y < height-1 && maze[y+1][x] == '#'
				left := x > 0 && maze[y][x-1] == '#'
				right := x < width-1 && maze[y][x+1] == '#'

				// Check if the wall is on the outer boundary
				topBoundary := y == 0
				bottomBoundary := y == height-1
				leftBoundary := x == 0
				rightBoundary := x == width-1

				// Handle outer walls with double-line pseudographics
				switch {
				// Handle tunnels
				case !topBoundary && !bottomBoundary && leftBoundary && !rightBoundary && !bottom:
					newRow[x] = '╜'
				case !topBoundary && !bottomBoundary && leftBoundary && !rightBoundary && !top:
					newRow[x] = '╖'
				case !topBoundary && !bottomBoundary && !leftBoundary && rightBoundary && !bottom:
					newRow[x] = '╙'
				case !topBoundary && !bottomBoundary && !leftBoundary && rightBoundary && !top:
					newRow[x] = '╓'

				case topBoundary && !bottomBoundary && leftBoundary && !rightBoundary:
					newRow[x] = '╔' // Top-left corner
				case topBoundary && !bottomBoundary && !leftBoundary && rightBoundary:
					newRow[x] = '╗' // Top-right corner
				case !topBoundary && bottomBoundary && leftBoundary && !rightBoundary:
					newRow[x] = '╚' // Bottom-left corner
				case !topBoundary && bottomBoundary && !leftBoundary && rightBoundary:
					newRow[x] = '╝' // Bottom-right corner
				case (topBoundary || bottomBoundary) && !leftBoundary && !rightBoundary && !top && !bottom:
					newRow[x] = '═' // Horizontal boundary
				case !topBoundary && !bottomBoundary && (leftBoundary || rightBoundary) && !left && !right && (top || bottom):
					newRow[x] = '║' // Vertical boundary

				// Handle connections between outer and inner walls
				case !topBoundary && !bottomBoundary && leftBoundary && !rightBoundary && right:
					newRow[x] = '╟' // Connects ║ with ─
				case topBoundary && !bottomBoundary && !leftBoundary && !rightBoundary && bottom:
					newRow[x] = '╤' // Connects ═ with │
				case !topBoundary && !bottomBoundary && !leftBoundary && rightBoundary && left:
					newRow[x] = '╢' // Connects ║ with ─
				case !topBoundary && bottomBoundary && !leftBoundary && !rightBoundary && top:
					newRow[x] = '╧' // Connects ═ with │

				// Handle standalone walls
				case !top && !bottom && !left && !right:
					newRow[x] = '─'

				// Handle inner walls
				default:
					switch {
					case !left && !right && (top || bottom):
						newRow[x] = '│'
					case !top && !bottom && (left || right):
						newRow[x] = '─'
					case !top && bottom && !left && right:
						newRow[x] = '┌'
					case !top && bottom && left && !right:
						newRow[x] = '┐'
					case top && !bottom && left && !right:
						newRow[x] = '┘'
					case top && !bottom && !left && right:
						newRow[x] = '└'
					case top && bottom && !left && right:
						newRow[x] = '├'
					case !top && bottom && left && right:
						newRow[x] = '┬'
					case top && bottom && left && !right:
						newRow[x] = '┤'
					case top && !bottom && left && right:
						newRow[x] = '┴'
					case top && bottom && left && right:
						newRow[x] = '┼'
					default:
						newRow[x] = '─'
					}
				}
			}
		}
		newMaze[y] = string(newRow)
	}

	return newMaze
}

func main() {
	// New example maze
	maze := []string{
		"###################",
		"#o.......#.......o#",
		"#.###.#..#..#.###.#",
		"#B.......#.......I#",
		"#.###.#..#..#.###.#",
		"  ....#..#..#....  ",
		"#C######.P.######.#",
		"#.....#..Y..#.....#",
		"#.###.#..#..#.###.#",
		"#o.......#.......o#",
		"###################",
	}

	p := tea.NewProgram(initialModel(maze))
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

// var pacman = []rune{'⭘', '◐', '◓', '◑', '◒'}

// // Ghosts
// var letters = []rune{'B', 'P', 'I', 'Y'}     // Letters
// var ghosts = []rune{'👿', '👽', '🤖', '👾'}      // Ghosts
// var hebrew = []rune{'ℵ', 'ℶ', 'ℷ', 'ℸ'}      // Hebrew ghosts
// var greek = []rune{'Α', 'Β', 'Γ', 'Δ'}       // Greek ghosts
// var control = []rune{'␊', '␋', '␌', '␍'}     // Control ghosts
// var currency = []rune{'$', '€', '£', '¥'}    // Currency ghosts
// var mathematics = []rune{'∀', '√', '∂', '∫'} // Math ghosts
