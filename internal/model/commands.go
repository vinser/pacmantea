package model

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinser/pacmantea/internal/sound"
	"github.com/vinser/pacmantea/internal/utils"
)

const (
	pacmanBlinkTickDuration = time.Second / 2
)

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
func (m *Model) pacmanBlinkTick() tea.Cmd {
	return tea.Tick(pacmanBlinkTickDuration, func(_ time.Time) tea.Msg {
		select {
		case <-m.Ctx.Done():
			return nil
		default:
			return pacmanBlinkMsg{}
		}
	})
}

// Message type for ghost movement
type ghostMoveMsg struct{}

// Command to trigger ghost movement ticks
func (m *Model) ghostMoveTick() tea.Cmd {
	if m.GameWin {
		// Do not schedule ghost movement if the game is won
		return func() tea.Msg { return nil }
	}
	return tea.Tick(time.Second/time.Duration(m.Difficulties[m.Levels[m.CurrentLevel].DifficultyName].GhostSpeed), func(_ time.Time) tea.Msg {
		select {
		case <-m.Ctx.Done():
			return nil
		default:
			return ghostMoveMsg{}
		}
	})
}

// Message types for rampant and cooldown states
type rampantEndMsg struct{}

// Command to start the rampant timer
func (m *Model) startRampantTimer() tea.Cmd {
	return tea.Tick(time.Duration(m.Difficulties[m.Levels[m.CurrentLevel].DifficultyName].RampantDuration)*time.Second, func(_ time.Time) tea.Msg {
		select {
		case <-m.Ctx.Done():
			return nil
		default:
			return rampantEndMsg{}
		}
	})
}

type cooldownEndMsg struct{}

// Command to start the cooldown timer
func (m *Model) startCooldownTimer() tea.Cmd {
	return tea.Tick(time.Duration(m.Difficulties[m.Levels[m.CurrentLevel].DifficultyName].CooldownDuration)*time.Second, func(_ time.Time) tea.Msg {
		select {
		case <-m.Ctx.Done():
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
func (m *Model) startGhostRevivalTimer(ghostName string, duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(_ time.Time) tea.Msg {
		select {
		case <-m.Ctx.Done():
			return nil
		default:
			return ghostReviveMsg{ghostName: ghostName}
		}
	})
}

// Check for collision with ghosts
func (m *Model) checkGhostCollisions() tea.Cmd {
	select {
	case <-m.Ctx.Done():
		return nil
	default:

		for name, g := range m.Ghosts {
			if g.Dead {
				continue
			}
			if m.Pacman.Position == g.Position {
				if m.Pacman.RampantState || m.Pacman.CooldownState {
					g.Dead = true
					ghostsEaten++
					m.LevelScore += ghostBonus * ghostsEaten
					go m.PlaySound(sound.EATGHOST)
					m.Ghosts[name] = g
					m.Maze[g.Position.Y] = utils.ReplaceAtIndex(m.Maze[g.Position.Y], ' ', g.Position.X) // Remove ghost from maze
					return m.startGhostRevivalTimer(name, time.Duration(m.Difficulties[m.Levels[m.CurrentLevel].DifficultyName].RevivalTimer)*time.Second)
				} else {
					m.GameOver = true
					go m.PlaySound(sound.DEATH)
				}
			}
		}
	}
	return nil
}
