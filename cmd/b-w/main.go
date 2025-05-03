package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type point struct {
	x, y int
}

type model struct {
	maze         []string // Maze as an array of strings (current state)
	originalMaze []string // Original maze (immutable)
	player       point    // Pac-Man's position
	dots         []point  // Positions of dots
	ghosts       []point  // Positions of ghosts
	score        int      // Score
	gameOver     bool     // Game over flag
	win          bool     // Win flag
}

func initialModel(maze []string) model {
	rand.Seed(time.Now().UnixNano())

	// Ensure the maze has a minimum size of 5x5
	if len(maze) < 5 || len(maze[0]) < 5 {
		panic("The maze must be at least 5x5")
	}

	var player point
	dots := []point{}
	ghosts := []point{}

	// Create a copy of the original maze
	originalMaze := make([]string, len(maze))
	copy(originalMaze, maze)

	for y, row := range maze {
		for x, char := range row {
			switch char {
			case 'C':
				player = point{x: x, y: y}
				// Replace 'C' with '.' after initialization
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'o':
				dots = append(dots, point{x: x, y: y})
			case 'G':
				ghosts = append(ghosts, point{x: x, y: y})
				// Replace 'G' with '.' after initialization
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case ' ':
				// Replace spaces with '.'
				maze[y] = replaceAtIndex(maze[y], '.', x)
			}
		}
	}

	// Convert maze walls to pseudographics for the current maze only
	maze = replaceWallsWithGraphics(maze)

	return model{
		maze:         maze,
		originalMaze: originalMaze, // Keep originalMaze unchanged
		player:       player,
		dots:         dots,
		ghosts:       ghosts,
		score:        0,
		gameOver:     false,
		win:          false,
	}
}

func (m model) Init() tea.Cmd {
	// Start the timer for ghost movement
	return ghostMoveTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver || m.win {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				// Restart the game using originalMaze
				newModel := initialModel(m.originalMaze)
				// Start the timer for ghost movement
				return newModel, ghostMoveTick()
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
			if m.canMove(m.player.x, m.player.y-1) {
				m.player.y--
			}
		case "down":
			if m.canMove(m.player.x, m.player.y+1) {
				m.player.y++
			}
		case "left":
			if m.canMove(m.player.x-1, m.player.y) {
				m.player.x--
			}
		case "right":
			if m.canMove(m.player.x+1, m.player.y) {
				m.player.x++
			}
		}

		// Check for dot collection
		for i := len(m.dots) - 1; i >= 0; i-- {
			if m.player == m.dots[i] {
				m.score++
				m.maze[m.player.y] = replaceAtIndex(m.maze[m.player.y], '.', m.player.x) // Replace dot with a space
				m.dots = append(m.dots[:i], m.dots[i+1:]...)
			}
		}

		// Check for win condition
		if len(m.dots) == 0 {
			m.win = true
		}

		// Check for collision with ghosts
		for _, ghost := range m.ghosts {
			if m.player == ghost {
				m.gameOver = true
			}
		}

	case ghostMoveMsg:
		// Move ghosts
		for i := range m.ghosts {
			directions := []point{
				{0, -1}, {0, 1}, {-1, 0}, {1, 0},
			}
			rand.Shuffle(len(directions), func(i, j int) {
				directions[i], directions[j] = directions[j], directions[i]
			})

			for _, dir := range directions {
				newX, newY := m.ghosts[i].x+dir.x, m.ghosts[i].y+dir.y
				if m.canMove(newX, newY) {
					m.ghosts[i] = point{x: newX, y: newY}
					break
				}
			}
		}

		// Check for collision with ghosts after their movement
		for _, ghost := range m.ghosts {
			if m.player == ghost {
				m.gameOver = true
			}
		}

		// Start the next tick for ghost movement
		return m, ghostMoveTick()
	}

	return m, nil
}

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
	for _, dot := range m.dots {
		grid[dot.y] = replaceAtIndex(grid[dot.y], 'o', dot.x)
	}

	// Place ghosts
	for _, ghost := range m.ghosts {
		grid[ghost.y] = replaceAtIndex(grid[ghost.y], 'G', ghost.x)
	}

	// Place the player
	grid[m.player.y] = replaceAtIndex(grid[m.player.y], 'C', m.player.x)

	// Build the string for display
	view := strings.Join(grid, "\n")
	view += fmt.Sprintf("\nScore: %d\n", m.score)
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
	if y < 0 || y >= len(m.maze) || x < 0 || x >= len([]rune(m.maze[0])) {
		return false // Out of bounds
	}

	// Convert the maze row to runes to handle multi-byte characters
	row := []rune(m.maze[y])
	char := row[x]

	// List of all wall characters (pseudographics)
	wallChars := []rune{
		'│', '─', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼',
		'║', '═', '╔', '╗', '╚', '╝', '╟', '╢', '╤', '╧',
	}

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

// Command to trigger ghost movement ticks
func ghostMoveTick() tea.Cmd {
	return tea.Tick(time.Second/2, func(_ time.Time) tea.Msg {
		return ghostMoveMsg{}
	})
}

// Convert maze walls to pseudographics
func replaceWallsWithGraphics(maze []string) []string {
	height := len(maze)
	width := len(maze[0])

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
				case !topBoundary && !bottomBoundary && (leftBoundary || rightBoundary) && !left && !right:
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
		"#G.......#.......G#",
		"#.###.#..#..#.###.#",
		"#o....#..#..#....o#",
		"#C######.G.######.#",
		"#.....#.....#.....#",
		"#.###.#..#..#.###.#",
		"#o.......#.......o#",
		"###################",
	}

	p := tea.NewProgram(initialModel(maze))
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
