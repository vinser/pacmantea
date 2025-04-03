package main

import (
	"context"
	"log"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed config.yml
var CONFIG_DATA []byte

// Simulate ghosts love
var ghostsLove bool = false

func newModel() model {
	var config Config
	err := yaml.Unmarshal(CONFIG_DATA, &config)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return initialModel(config, 0)

}

func createPlayerAt(pos point) player {
	return player{
		entity: entity{
			position: pos,
			style:    playerStyle,
			name:     "Pac-Man",
			badge:    'C',
		},
		blinkState: false,
	}
}

func createDotAt(pos point) dot {
	return dot{
		entity: entity{
			position: pos,
			style:    dotStyle,
			name:     "Dot",
			badge:    '.',
		},
	}
}

func createEnergizer(x, y int) energizer {
	return energizer{
		entity: entity{
			position: point{x: x, y: y},
			style:    energyStyle,
			name:     "Energizer",
			badge:    'o',
		},
	}
}

func createGhostAt(pos point, style lipgloss.Style, name string, badge rune) ghost {
	return ghost{
		entity: entity{
			position: pos,
			style:    style,
			name:     name,
			badge:    badge,
		},
		revivalPoint: pos,
	}
}

// Update the initialModel function
func initialModel(config Config, currntLevel int) model {
	maze := make([]string, len(config.Levels[currntLevel].Maze))
	copy(maze, config.Levels[currntLevel].Maze)
	// Ensure the maze has a minimum size of 5x5
	if len(maze) < 5 || len(maze[0]) < 5 {
		log.Fatal("The maze must be at least 5x5")
	}

	var playerEntity player
	dots := []dot{}
	energizers := []energizer{}
	ghosts := make(map[string]ghost)

	playerPlaced := false
	ghostsPlaced := map[string]bool{"Blinky": false, "Inky": false, "Pinky": false, "Clyde": false}

	for y, row := range maze {
		for x, char := range row {
			switch char {
			case 'C':
				playerEntity = createPlayerAt(point{x: x, y: y})
				playerPlaced = true
				dots = append(dots, createDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case '.':
				dots = append(dots, createDotAt(point{x: x, y: y}))
			case 'o':
				energizers = append(energizers, createEnergizer(x, y))
			case 'B':
				ghosts["Blinky"] = createGhostAt(point{x: x, y: y}, blinkyStyle, "Blinky", 'B')
				ghostsPlaced["Blinky"] = true
				dots = append(dots, createDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'I':
				ghosts["Inky"] = createGhostAt(point{x: x, y: y}, inkyStyle, "Inky", 'I')
				ghostsPlaced["Inky"] = true
				dots = append(dots, createDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'P':
				ghosts["Pinky"] = createGhostAt(point{x: x, y: y}, pinkyStyle, "Pinky", 'P')
				ghostsPlaced["Pinky"] = true
				dots = append(dots, createDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'Y':
				ghosts["Clyde"] = createGhostAt(point{x: x, y: y}, clydeStyle, "Clyde", 'Y')
				ghostsPlaced["Clyde"] = true
				dots = append(dots, createDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			}
		}
	}

	// Randomly place the player from edges to center if not placed
	if !playerPlaced {
		playerEntity = placePlayerRandomly(maze)
	}

	// Randomly place ghosts near the center if not placed
	for name, placed := range ghostsPlaced {
		if !placed {
			var style lipgloss.Style
			var badge rune
			switch name {
			case "Blinky":
				style, badge = blinkyStyle, 'B'
			case "Inky":
				style, badge = inkyStyle, 'I'
			case "Pinky":
				style, badge = pinkyStyle, 'P'
			case "Clyde":
				style, badge = clydeStyle, 'Y'
			}
			ghosts[name] = placeGhostRandomly(maze, name, style, badge)
		}
	}

	// Convert maze walls to pseudographics for the current maze only
	maze = replaceWallsWithPseudographics(maze)
	ctx, cancel := context.WithCancel(context.Background())
	return model{
		ctx:         ctx,
		cancel:      cancel,
		Config:      config,
		currntLevel: currntLevel,
		maze:        maze,
		maxScore:    len(dots),
		player:      playerEntity,
		dots:        dots,
		energizers:  energizers,
		ghosts:      ghosts,
		score:       0,
		gameOver:    false,
		win:         false,
		lives:       5, // Initialize with 5 lives
	}
}

func placePlayerRandomly(maze []string) player {
	pos := point{}
	free := traverseOrder(maze, MazePerifery)
	if len(free) > 0 {
		pos = free[rng.Intn(min(4, len(free)))]
	}
	return createPlayerAt(pos)

}

func placeGhostRandomly(maze []string, name string, style lipgloss.Style, badge rune) ghost {
	pos := point{}
	free := traverseOrder(maze, MazeCenter)
	if len(free) > 0 {
		pos = free[rng.Intn(min(6, len(free)))]
	}
	return createGhostAt(pos, style, name, badge)

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
