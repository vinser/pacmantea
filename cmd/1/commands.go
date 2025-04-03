package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const playerBlinkTickDuration = time.Second / 2

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.playerBlinkTick(), // Start the timer for Pac-Man blinking
		m.ghostMoveTick(),   // Start the timer for ghost movement
	}
	return tea.Batch(cmds...)
}

// Message type for blinking
type playerBlinkMsg struct{}

// Command to trigger Pac-Man blinking
func (m *model) playerBlinkTick() tea.Cmd {
	return tea.Tick(playerBlinkTickDuration, func(_ time.Time) tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		default:
			return playerBlinkMsg{}
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
	return tea.Tick(time.Duration(float64(time.Second)/float64(m.Difficulties[m.Levels[m.currntLevel].Difficulty].GhostSpeed)), func(_ time.Time) tea.Msg {
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
	return tea.Tick(time.Duration(m.Difficulties[m.Levels[m.currntLevel].Difficulty].RampantDuration)*time.Second, func(_ time.Time) tea.Msg {
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
	return tea.Tick(time.Duration(m.Difficulties[m.Levels[m.currntLevel].Difficulty].CooldownDuration)*time.Second, func(_ time.Time) tea.Msg {
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
			if m.player.position == g.position {
				if m.player.rampantState || m.player.cooldownState {
					g.dead = true
					m.ghosts[name] = g
					m.maze[g.position.y] = replaceAtIndex(m.maze[g.position.y], ' ', g.position.x) // Remove ghost from maze
					return m.startGhostRevivalTimer(name, time.Duration(m.Difficulties[m.Levels[m.currntLevel].Difficulty].RevivalTimer)*time.Second)
				} else {
					m.gameOver = true
				}
			}
		}
	}
	return nil
}
