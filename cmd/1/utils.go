package main

import (
	"math/rand"
	"sort"
	"time"
)

type proximity bool

const (
	MazePerifery proximity = true
	MazeCenter   proximity = false
)

// Sort available (free) positions based on their distance from the center of the maze in a specified proximity.
func traverseOrder(maze []string, prox proximity) []point {
	rows := len(maze)
	cols := len(maze[0])

	// Calculate the center of the rectangle
	centerPoint := point{x: int(float64(cols-1) / 2.0), y: int(float64(rows-1) / 2.0)}

	// Create a list of all cells with their coordinates and values
	var points []point
	for y, row := range maze {
		for x, r := range row {
			if r != '#' && r != 'o' && r != 'C' && r != 'B' && r != 'I' && r != 'P' && r != 'Y' { // Skip occupied points in the maze
				points = append(points, point{x: x, y: y})
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

// Sort directions by distance to the player
func sortDirectionsByDistance(p1, p2 point) []point {
	directions := []point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	sort.Slice(directions, func(i, j int) bool {
		d1 := distanceSquare(point{x: p1.x + directions[i].x, y: p1.y + directions[i].y}, p2)
		d2 := distanceSquare(point{x: p1.x + directions[j].x, y: p1.y + directions[j].y}, p2)
		return d1 < d2
	})
	return directions
}

func distanceSquare(p1, p2 point) float64 {
	return float64((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y))
}

// Seed local random number generator
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomDirections() []point {
	directions := []point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	rng.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	return directions
}

// Replace a character in a string at a given index
func replaceAtIndex(s string, r rune, index int) string {
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
