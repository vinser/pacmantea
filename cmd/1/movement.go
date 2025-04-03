package main

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

func (m model) tunnelMove(newX int) int {
	width := len([]rune(m.maze[0]))
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
	'╖': true, '╓': true, '╜': true, '╙': true, //Tunnels corners
}

// Check if movement is possible
func (m model) canMove(x, y int) bool {
	row := []rune(m.maze[y])
	char := row[x]
	return !wallChars[char] // Not a wall, movement is allowed
}
