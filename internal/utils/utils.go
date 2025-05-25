package utils

import (
	"math/rand"
	"sort"
	"time"
)

type Point struct {
	X, Y int
}

type Direction struct {
	X, Y int
}

type proximity bool

const (
	MazePerifery proximity = true
	MazeCenter   proximity = false
)

// Sort available (free) positions based on their distance from the center of the maze in a specified proximity.
func TraverseOrder(maze []string, prox proximity) []Point {
	rows := len(maze)
	cols := len(maze[0])

	// Calculate the center of the rectangle
	centerPoint := Point{X: int(float64(cols-1) / 2.0), Y: int(float64(rows-1) / 2.0)}

	// Create a list of all cells with their coordinates and values
	var points []Point
	for y, row := range maze {
		for x, r := range row {
			if r != '#' && r != 'o' && r != 'C' && r != 'B' && r != 'I' && r != 'P' && r != 'Y' { // Skip occupied points in the maze
				points = append(points, Point{X: x, Y: y})
			}
		}
	}
	// Sort points by distance from the center
	sort.Slice(points, func(i, j int) bool {
		distI := distanceSquare(points[i], centerPoint)
		distJ := distanceSquare(points[j], centerPoint)
		if prox {
			return distI > distJ // At periphery
		}
		return distI < distJ // Near center
	})

	return points
}

// Sort directions by distance to the pacman
func SortDirectionsByDistance(p1, p2 Point) []Point {
	directions := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	sort.Slice(directions, func(i, j int) bool {
		d1 := distanceSquare(Point{X: p1.X + directions[i].X, Y: p1.Y + directions[i].Y}, p2)
		d2 := distanceSquare(Point{X: p1.X + directions[j].X, Y: p1.Y + directions[j].Y}, p2)
		return d1 < d2
	})
	return directions
}

func distanceSquare(p1, p2 Point) float64 {
	return float64((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
}

// Seed local random number generator
var Rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomDirections() []Point {
	directions := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	Rng.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	return directions
}

// Replace a character in a string at a given index
func ReplaceAtIndex(s string, r rune, index int) string {
	runes := []rune(s)
	if index < 0 || index >= len(runes) {
		return s // Return the original string if the index is out of bounds
	}
	runes[index] = r
	return string(runes)
}

// Compute the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Compute the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
