package model

import (
	"context"
	"log"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/vinser/pacmantea/internal/config"
	"github.com/vinser/pacmantea/internal/sound"
	"github.com/vinser/pacmantea/internal/state"
	"github.com/vinser/pacmantea/internal/ui"
	"github.com/vinser/pacmantea/internal/utils"
)

type Entity struct {
	Position utils.Point
	Move     utils.Direction
	Style    lipgloss.Style
	Name     string
	Badge    rune
}

type Pacman struct {
	Entity
	ChewState     bool
	RampantState  bool
	CooldownState bool
}

type Ghost struct {
	Entity
	Dead         bool
	RevivalPoint utils.Point
}

type Dot struct {
	Entity
}

type Energizer struct {
	Entity
}

type Model struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	config.Config
	state.State
	CurrentLevel int
	CurrentSart  time.Time
	Maze         []string
	Pacman       Pacman
	Dots         []Dot
	Energizers   []Energizer
	Ghosts       map[string]Ghost
	LevelScore   int
	GameScore    int
	GameOver     bool
	LevelWin     bool
	GameWin      bool
	Lives        int
	Sounds       map[string]sound.Sound
}

func New() *Model {
	state := state.Load()
	config := config.Load()
	// Initialize the game model with the loaded configuration and saved game
	return InitialModel(config, state)
}

// InitialModel returns the initial model for the game
func InitialModel(config config.Config, state state.State) *Model {
	currntLevel := 0
	for i, level := range config.Levels {
		if level.Name == state.LevelName {
			currntLevel = i
			break
		}
	}
	if state.LevelName == "" {
		state.LevelName = config.Levels[currntLevel].Name
	}
	maze := make([]string, len(config.Levels[currntLevel].Maze))
	copy(maze, config.Levels[currntLevel].Maze)
	// Ensure the maze has a minimum size of 5x5
	if len(maze) < 5 || len(maze[0]) < 5 {
		log.Fatal("The maze must be at least 5x5")
	}

	var pacmanEntity Pacman
	dots := []Dot{}
	energizers := []Energizer{}
	ghosts := make(map[string]Ghost)

	pacmanPlaced := false
	ghostsPlaced := map[string]bool{"Blinky": false, "Inky": false, "Pinky": false, "Clyde": false}

	for y, row := range maze {
		for x, char := range row {
			switch char {
			case 'C':
				pacmanEntity = initPacmanAt(utils.Point{X: x, Y: y})
				pacmanPlaced = true
				dots = append(dots, initDotAt(utils.Point{X: x, Y: y}))
				maze[y] = utils.ReplaceAtIndex(maze[y], '.', x)
			case '.':
				dots = append(dots, initDotAt(utils.Point{X: x, Y: y}))
			case 'o':
				energizers = append(energizers, initEnergizerAt(x, y))
			case 'B':
				ghosts["Blinky"] = initGhostAt(utils.Point{X: x, Y: y}, ui.BlinkyStyle, "Blinky", 'B')
				ghostsPlaced["Blinky"] = true
				dots = append(dots, initDotAt(utils.Point{X: x, Y: y}))
				maze[y] = utils.ReplaceAtIndex(maze[y], '.', x)
			case 'I':
				ghosts["Inky"] = initGhostAt(utils.Point{X: x, Y: y}, ui.InkyStyle, "Inky", 'I')
				ghostsPlaced["Inky"] = true
				dots = append(dots, initDotAt(utils.Point{X: x, Y: y}))
				maze[y] = utils.ReplaceAtIndex(maze[y], '.', x)
			case 'P':
				ghosts["Pinky"] = initGhostAt(utils.Point{X: x, Y: y}, ui.PinkyStyle, "Pinky", 'P')
				ghostsPlaced["Pinky"] = true
				dots = append(dots, initDotAt(utils.Point{X: x, Y: y}))
				maze[y] = utils.ReplaceAtIndex(maze[y], '.', x)
			case 'Y':
				ghosts["Clyde"] = initGhostAt(utils.Point{X: x, Y: y}, ui.ClydeStyle, "Clyde", 'Y')
				ghostsPlaced["Clyde"] = true
				dots = append(dots, initDotAt(utils.Point{X: x, Y: y}))
				maze[y] = utils.ReplaceAtIndex(maze[y], '.', x)
			}
		}
	}

	// Randomly place the pacman from edges to center if not placed
	if !pacmanPlaced {
		pacmanEntity = placePacmanRandomly(maze)
	}

	// Randomly place ghosts near the center if not placed
	for name, placed := range ghostsPlaced {
		if !placed {
			var style lipgloss.Style
			var badge rune
			switch name {
			case "Blinky":
				style, badge = ui.BlinkyStyle, 'B'
			case "Inky":
				style, badge = ui.InkyStyle, 'I'
			case "Pinky":
				style, badge = ui.PinkyStyle, 'P'
			case "Clyde":
				style, badge = ui.ClydeStyle, 'Y'
			}
			ghosts[name] = placeGhostRandomly(maze, name, style, badge)
		}
	}

	// Convert maze walls to pseudographics for the current maze only
	maze = replaceWallsWithPseudographics(maze)
	sounds, err := sound.LoadSamples()

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Model{
		Ctx:          ctx,
		Cancel:       cancel,
		Config:       config,
		State:        state,
		CurrentLevel: currntLevel,
		CurrentSart:  time.Now(),
		Maze:         maze,
		Pacman:       pacmanEntity,
		Dots:         dots,
		Energizers:   energizers,
		Ghosts:       ghosts,
		LevelScore:   0,
		GameOver:     false,
		LevelWin:     false,
		Lives:        5, // Initialize with 5 lives
		Sounds:       sounds,
	}
}

func initPacmanAt(pos utils.Point) Pacman {
	return Pacman{
		Entity: Entity{
			Position: pos,
			Style:    ui.PacmanStyle,
			Name:     "Pac-Man",
			Badge:    'C',
		},
		ChewState: false,
	}
}

func initDotAt(pos utils.Point) Dot {
	return Dot{
		Entity: Entity{
			Position: pos,
			Style:    ui.DotStyle,
			Name:     "Dot",
			Badge:    '.',
		},
	}
}

func initEnergizerAt(x, y int) Energizer {
	return Energizer{
		Entity: Entity{
			Position: utils.Point{X: x, Y: y},
			Style:    ui.EnergyStyle,
			Name:     "Energizer",
			Badge:    'o',
		},
	}
}

func initGhostAt(pos utils.Point, style lipgloss.Style, name string, badge rune) Ghost {
	return Ghost{
		Entity: Entity{
			Position: pos,
			Style:    style,
			Name:     name,
			Badge:    badge,
		},
		RevivalPoint: pos,
	}
}

func placePacmanRandomly(maze []string) Pacman {
	pos := utils.Point{}
	free := utils.TraverseOrder(maze, utils.MazePerifery)
	if len(free) > 0 {
		pos = free[utils.Rng.Intn(min(4, len(free)))]
	}
	return initPacmanAt(pos)

}

func placeGhostRandomly(maze []string, name string, style lipgloss.Style, badge rune) Ghost {
	pos := utils.Point{}
	free := utils.TraverseOrder(maze, utils.MazeCenter)
	if len(free) > 0 {
		pos = free[utils.Rng.Intn(min(6, len(free)))]
	}
	return initGhostAt(pos, style, name, badge)

}

// Convert maze walls to pseudographics
func replaceWallsWithPseudographics(maze []string) []string {
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
				case !topBoundary && !bottomBoundary && (leftBoundary || rightBoundary) && (left || right) && !bottom:
					newRow[x] = '╨'
				case !topBoundary && !bottomBoundary && (leftBoundary || rightBoundary) && (left || right) && !top:
					newRow[x] = '╥'

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
