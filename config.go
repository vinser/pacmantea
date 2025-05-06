package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed config.yml
var CONFIG_DATA []byte

var configPath string = path.Join("config", "config.yml")

// Simulate ghosts love
var ghostsLove bool = false

func writeDefaultConfig() error {
	// Write the embedded CONFIG_DATA to a file
	if err := os.MkdirAll("config", os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(CONFIG_DATA)
	return err
}

func newModel() model {
	var config Config

	// Determine the config path based on the environment variable
	configPath := os.Getenv("PACMANTEA_CONFIG_PATH")
	if configPath == "" {
		configPath = path.Join("config", "config.yml")
	}

	// Load the configuration from the file or use the embedded one as a fallback
	var data []byte
	var err error
	if _, err = os.Stat(configPath); err == nil {
		data, err = os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", configPath, err)
		}
	} else {
		data = CONFIG_DATA
	}

	// Unmarshal the configuration
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse %s: %v", configPath, err)
	}

	// Initialize the game model with the loaded configuration
	return initialModel(config, 0)
}

func initPacmanAt(pos point) pacman {
	return pacman{
		entity: entity{
			position: pos,
			style:    pacmanStyle,
			name:     "Pac-Man",
			badge:    'C',
		},
		chewState: false,
	}
}

func initDotAt(pos point) dot {
	return dot{
		entity: entity{
			position: pos,
			style:    dotStyle,
			name:     "Dot",
			badge:    '.',
		},
	}
}

func initEnergizerAt(x, y int) energizer {
	return energizer{
		entity: entity{
			position: point{x: x, y: y},
			style:    energyStyle,
			name:     "Energizer",
			badge:    'o',
		},
	}
}

func initGhostAt(pos point, style lipgloss.Style, name string, badge rune) ghost {
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

	var pacmanEntity pacman
	dots := []dot{}
	energizers := []energizer{}
	ghosts := make(map[string]ghost)

	pacmanPlaced := false
	ghostsPlaced := map[string]bool{"Blinky": false, "Inky": false, "Pinky": false, "Clyde": false}

	for y, row := range maze {
		for x, char := range row {
			switch char {
			case 'C':
				pacmanEntity = initPacmanAt(point{x: x, y: y})
				pacmanPlaced = true
				dots = append(dots, initDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case '.':
				dots = append(dots, initDotAt(point{x: x, y: y}))
			case 'o':
				energizers = append(energizers, initEnergizerAt(x, y))
			case 'B':
				ghosts["Blinky"] = initGhostAt(point{x: x, y: y}, blinkyStyle, "Blinky", 'B')
				ghostsPlaced["Blinky"] = true
				dots = append(dots, initDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'I':
				ghosts["Inky"] = initGhostAt(point{x: x, y: y}, inkyStyle, "Inky", 'I')
				ghostsPlaced["Inky"] = true
				dots = append(dots, initDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'P':
				ghosts["Pinky"] = initGhostAt(point{x: x, y: y}, pinkyStyle, "Pinky", 'P')
				ghostsPlaced["Pinky"] = true
				dots = append(dots, initDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
			case 'Y':
				ghosts["Clyde"] = initGhostAt(point{x: x, y: y}, clydeStyle, "Clyde", 'Y')
				ghostsPlaced["Clyde"] = true
				dots = append(dots, initDotAt(point{x: x, y: y}))
				maze[y] = replaceAtIndex(maze[y], '.', x)
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
	sounds, err := loadSoundSamples()

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return model{
		ctx:          ctx,
		cancel:       cancel,
		Config:       config,
		currentLevel: currntLevel,
		maze:         maze,
		maxScore:     len(dots),
		pacman:       pacmanEntity,
		dots:         dots,
		energizers:   energizers,
		ghosts:       ghosts,
		score:        0,
		gameOver:     false,
		win:          false,
		lives:        5, // Initialize with 5 lives
		sounds:       sounds,
	}
}

func placePacmanRandomly(maze []string) pacman {
	pos := point{}
	free := traverseOrder(maze, MazePerifery)
	if len(free) > 0 {
		pos = free[rng.Intn(min(4, len(free)))]
	}
	return initPacmanAt(pos)

}

func placeGhostRandomly(maze []string, name string, style lipgloss.Style, badge rune) ghost {
	pos := point{}
	free := traverseOrder(maze, MazeCenter)
	if len(free) > 0 {
		pos = free[rng.Intn(min(6, len(free)))]
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
