package ui

import "github.com/charmbracelet/lipgloss"

// Define styles for different elements
var (
	WallStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))            // Green
	PacmanStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true) // Yellow
	DotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))           // White
	EnergyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true) // Blue
)

// Define styles for different ghosts
var (
	BlinkyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)   // Red
	InkyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)   // Cyan
	PinkyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Bold(true) // Pink
	ClydeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Orange
)
