package main

func (m model) chaosMove(p point) point {
	directions := randomDirections()
	return m.ghostMove(p, directions)
}

func (m model) straitMove(p point) point {
	destination := m.player.position
	directions := sortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m model) predictMove(p point) point {
	destination := point{x: m.player.position.x + m.player.move.x, y: m.player.position.y + m.player.move.y}
	directions := sortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m model) cagyMove(p point) point {
	destination := point{x: m.player.position.x - 2*m.player.move.x, y: m.player.position.y - 2*m.player.move.y}
	directions := sortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m model) escapeMove(p point) point {
	destination := point{x: 2*p.x - m.player.position.x, y: 2*p.y - m.player.position.y}
	directions := sortDirectionsByDistance(p, destination)
	return m.ghostMove(p, directions)
}

func (m model) ghostMove(from point, directions []point) point {
	for _, dir := range directions {
		to := point{x: from.x + dir.x, y: from.y + dir.y}
		to.x = m.tunnelMove(to.x)
		if m.isGhostHere(to) {
			continue
		}
		if m.canMove(to.x, to.y) {
			return to
		}
	}
	return from
}

func (m model) isGhostHere(p point) bool {
	for _, g := range m.ghosts {
		if g.position.x == p.x && g.position.y == p.y {
			return true
		}
	}
	return false
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
	'╖': true, '╓': true, '╜': true, '╙': true, '╨': true, '╥': true, //Tunnels corners
}

// Check if movement is possible
func (m model) canMove(x, y int) bool {
	row := []rune(m.maze[y])
	char := row[x]
	return !wallChars[char] // Not a wall, movement is allowed
}
