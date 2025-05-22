package main

import (
	"context"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Level struct {
	Name           string   `yaml:"name"`
	DifficultyName string   `yaml:"difficulty"`
	Maze           []string `yaml:"maze"`
	PacmanBadge    string   `yaml:"pacman_badge"` // Badge style for Pac-Man
	GhostBadges    string   `yaml:"ghost_badges"` // Badge style for ghosts
}

type Config struct {
	Badges       Badges                `yaml:"badges"`
	Difficulties map[string]Difficulty `yaml:"difficulties"`
	Levels       []Level               `yaml:"levels"`
}

type Difficulty struct {
	GhostSpeed       int `yaml:"ghost_speed"`
	RampantDuration  int `yaml:"rampant_duration"`
	CooldownDuration int `yaml:"cooldown_duration"`
	RevivalTimer     int `yaml:"revival_timer"`
	SpeedBonus       int `yaml:"speed_bonus"` // base points for each second in formula speedBonus * wonGames * seconds
}

type Badges struct {
	Pacman map[string]map[string]string `yaml:"pacman"` // Badge styles for Pac-Man
	Ghosts map[string]map[string]string `yaml:"ghosts"` // Badge styles for ghosts
}

type point struct {
	x, y int
}

type direction struct {
	x, y int
}

type entity struct {
	position point
	move     direction
	style    lipgloss.Style
	name     string
	badge    rune
}

type pacman struct {
	entity
	chewState     bool
	rampantState  bool
	cooldownState bool
}

type ghost struct {
	entity
	dead         bool
	revivalPoint point
}

type dot struct {
	entity
}

type energizer struct {
	entity
}

type model struct {
	ctx    context.Context
	cancel context.CancelFunc
	Config
	State
	currentLevel int
	currentSart  time.Time
	maze         []string
	pacman       pacman
	dots         []dot
	energizers   []energizer
	ghosts       map[string]ghost
	levelScore   int
	gameScore    int
	gameOver     bool
	win          bool
	winGame      bool
	lives        int
	sounds       map[string]sound
}
