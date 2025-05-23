package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	pacmanBlinkTickDuration = time.Second / 2
)

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.pacmanBlinkTick(), // Start the timer for Pac-Man blinking
		m.ghostMoveTick(),   // Start the timer for ghost movement
		splashScreen(),
	}
	return tea.Batch(cmds...)
}

// Message type for starting the game
type startGameMsg struct{}

// Command to start the game
func splashScreen() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
		return startGameMsg{}
	})
}

// Message type for blinking
type pacmanBlinkMsg struct{}

// Command to trigger Pac-Man blinking
func (m *model) pacmanBlinkTick() tea.Cmd {
	return tea.Tick(pacmanBlinkTickDuration, func(_ time.Time) tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		default:
			return pacmanBlinkMsg{}
		}
	})
}

// Message type for ghost movement
type ghostMoveMsg struct{}

// Command to trigger ghost movement ticks
func (m *model) ghostMoveTick() tea.Cmd {
	if m.winGame {
		// Do not schedule ghost movement if the game is won
		return func() tea.Msg { return nil }
	}
	return tea.Tick(time.Second/time.Duration(m.Difficulties[m.Levels[m.currentLevel].DifficultyName].GhostSpeed), func(_ time.Time) tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		default:
			return ghostMoveMsg{}
		}
	})
}

// Message types for rampant and cooldown states
type rampantEndMsg struct{}

// Command to start the rampant timer
func (m *model) startRampantTimer() tea.Cmd {
	return tea.Tick(time.Duration(m.Difficulties[m.Levels[m.currentLevel].DifficultyName].RampantDuration)*time.Second, func(_ time.Time) tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		default:
			return rampantEndMsg{}
		}
	})
}

type cooldownEndMsg struct{}

// Command to start the cooldown timer
func (m *model) startCooldownTimer() tea.Cmd {
	return tea.Tick(time.Duration(m.Difficulties[m.Levels[m.currentLevel].DifficultyName].CooldownDuration)*time.Second, func(_ time.Time) tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		default:
			return cooldownEndMsg{}
		}
	})
}

// Message type for ghost revival
type ghostReviveMsg struct {
	ghostName string
}

// Command to start the ghost revival timer
func (m *model) startGhostRevivalTimer(ghostName string, duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(_ time.Time) tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		default:
			return ghostReviveMsg{ghostName: ghostName}
		}
	})
}

// Check for collision with ghosts
func (m *model) checkGhostCollisions() tea.Cmd {
	select {
	case <-m.ctx.Done():
		return nil
	default:

		for name, g := range m.ghosts {
			if g.dead {
				continue
			}
			if m.pacman.position == g.position {
				if m.pacman.rampantState || m.pacman.cooldownState {
					g.dead = true
					ghostsEaten++
					m.levelScore += ghostBonus * ghostsEaten
					go m.playSound(SOUND_EATGHOST)
					m.ghosts[name] = g
					m.maze[g.position.y] = replaceAtIndex(m.maze[g.position.y], ' ', g.position.x) // Remove ghost from maze
					return m.startGhostRevivalTimer(name, time.Duration(m.Difficulties[m.Levels[m.currentLevel].DifficultyName].RevivalTimer)*time.Second)
				} else {
					m.gameOver = true
					go m.playSound(SOUND_DEATH)
				}
			}
		}
	}
	return nil
}
