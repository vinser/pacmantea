package main

import (
	"context"

	"github.com/charmbracelet/lipgloss"
)

type Level struct {
	Name           string   `yaml:"name"`
	DifficultyName string   `yaml:"difficulty"`
	Maze           []string `yaml:"maze"`
}

type Config struct {
	Difficulties map[string]Difficulty `yaml:"difficulties"`
	Levels       []Level               `yaml:"levels"`
}

type Difficulty struct {
	GhostSpeed       int `yaml:"ghost_speed"`
	RampantDuration  int `yaml:"rampant_duration"`
	CooldownDuration int `yaml:"cooldown_duration"`
	RevivalTimer     int `yaml:"revival_timer"`
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

type player struct {
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
	currentLevel int
	maze         []string
	maxScore     int
	player       player
	dots         []dot
	energizers   []energizer
	ghosts       map[string]ghost
	score        int
	gameOver     bool
	win          bool
	winGame      bool
	lives        int
	sounds       map[string]sound
}
