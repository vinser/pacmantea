package model

import "github.com/vinser/pacmantea/internal/utils"

func (m Model) chaosMove(p utils.Point) utils.Point {
	directions := utils.RandomDirections()
	return m.ghostMove(p, directions)
}

func (m Model) straitMove(p utils.Point) utils.Point {
	destination := m.Pacman.Position
	directions := utils.SortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m Model) predictMove(p utils.Point) utils.Point {
	destination := utils.Point{X: m.Pacman.Position.X + m.Pacman.Move.X, Y: m.Pacman.Position.Y + m.Pacman.Move.Y}
	directions := utils.SortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m Model) cagyMove(p utils.Point) utils.Point {
	destination := utils.Point{X: m.Pacman.Position.X - 2*m.Pacman.Move.X, Y: m.Pacman.Position.Y - 2*m.Pacman.Move.Y}
	directions := utils.SortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m Model) escapeMove(p utils.Point) utils.Point {
	destination := utils.Point{X: 2*p.X - m.Pacman.Position.X, Y: 2*p.Y - m.Pacman.Position.Y}
	directions := utils.SortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m Model) ghostMove(from utils.Point, directions []utils.Point) utils.Point {
	for _, dir := range directions {
		to := utils.Point{X: from.X + dir.X, Y: from.Y + dir.Y}
		to.X = m.tunnelMove(to.X)
		if m.isGhostHere(to) {
			continue
		}
		if m.canMove(to.X, to.Y) {
			return to
		}
	}
	return from
}

func (m Model) isGhostHere(p utils.Point) bool {
	for _, g := range m.Ghosts {
		if g.Position.X == p.X && g.Position.Y == p.Y {
			return true
		}
	}
	return false
}

func (m Model) tunnelMove(newX int) int {
	width := len([]rune(m.Maze[0]))
	if newX < 0 {
		return width - 1
	}
	if newX >= width {
		return 0
	}
	return newX
}

// List of all wall characters (pseudographics)
var wallChars = map[rune]bool{
	'│': true, '─': true, '┌': true, '┐': true, '└': true, '┘': true, '├': true, '┤': true, '┬': true, '┴': true, '┼': true, // Inner wals
	'║': true, '═': true, '╔': true, '╗': true, '╚': true, '╝': true, '╟': true, '╢': true, '╤': true, '╧': true, // Outer wals
	'╖': true, '╓': true, '╜': true, '╙': true, '╨': true, '╥': true, //Tunnels corners
}

// Check if movement is possible
func (m Model) canMove(x, y int) bool {
	row := []rune(m.Maze[y])
	char := row[x]
	return !wallChars[char] // Not a wall, movement is allowed
}
