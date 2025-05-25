package model

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinser/pacmantea/internal/sound"
	"github.com/vinser/pacmantea/internal/state"
	"github.com/vinser/pacmantea/internal/utils"
)

const (
	ghostBonus = 40 // 40 points for first eated ghost in rampant state, 80 for second...
)

var ghostsEaten int = 0

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.GameOver {
		if m.Lives > 1 {
			// Wait for spacebar to restart the current level
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					lives := m.Lives - 1
					m.Cancel()                                  // Deduct a life
					newModel := InitialModel(m.Config, m.State) // Restart current level
					newModel.Lives = lives                      // Preserve remaining lives
					return newModel, newModel.Init()
				case "q", "ctrl+c":
					m.LevelName = m.Levels[m.CurrentLevel].Name
					return m, tea.Quit
				case "m":
					m.Mute = !m.Mute
				}
			}
			return m, nil
		}
		// No lives left, offer to restart the game
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				m.Cancel()
				m.LevelName = m.Levels[0].Name
				newModel := InitialModel(m.Config, m.State)
				return newModel, newModel.Init()
			case "q", "ctrl+c":
				return m, tea.Quit
			case "m":
				m.Mute = !m.Mute
			}
		}
		return m, nil
	}
	if m.GameWin {
		state.Save(m.State)
		m.LevelName = ""
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				m.Cancel()
				newModel := InitialModel(m.Config, m.State)
				newModel.GameWin = false // Reset the winGame flag
				// Start the timer for ghost movement and blinking
				return newModel, newModel.Init()
			case "q", "ctrl+c":
				return m, tea.Quit
			case "m":
				m.Mute = !m.Mute
			}
		}
		// Stop scheduling commands when the game is won
		return m, nil
	}
	if m.LevelWin {
		m.recordLevelElapsedTime()
		if m.CurrentLevel >= len(m.Levels)-1 {
			m.GameWin = true
			return m, nil
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				if m.LevelWin {
					if m.CurrentLevel < len(m.Levels)-1 {
						m.CurrentLevel++
						m.LevelName = m.Levels[m.CurrentLevel].Name
					}
					lives := m.Lives
					m.Cancel()
					newModel := InitialModel(m.Config, m.State)
					newModel.Lives = lives
					// Start the timer for ghost movement and blinking
					return newModel, newModel.Init()
				}
			case "q", "ctrl+c":
				if m.CurrentLevel < len(m.Levels)-1 {
					m.CurrentLevel++
					m.LevelName = m.Levels[m.CurrentLevel].Name
				}
				state.Save(m.State)
				return m, tea.Quit
			case "m":
				m.Mute = !m.Mute
			}
		}
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		sound.ClearSpeaker()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "m":
			m.Mute = !m.Mute
		case "up":
			if m.canMove(m.Pacman.Position.X, m.Pacman.Position.Y-1) {
				m.Pacman.Position.Y--
				m.Pacman.Move = utils.Direction{X: 0, Y: -1}
			}
		case "down":
			if m.canMove(m.Pacman.Position.X, m.Pacman.Position.Y+1) {
				m.Pacman.Position.Y++
				m.Pacman.Move = utils.Direction{X: 0, Y: 1}
			}
		case "left":
			newX := m.tunnelMove(m.Pacman.Position.X - 1)
			if m.canMove(newX, m.Pacman.Position.Y) {
				m.Pacman.Position.X = newX
				m.Pacman.Move = utils.Direction{X: -1, Y: 0}
			}
		case "right":
			newX := m.tunnelMove(m.Pacman.Position.X + 1)
			if m.canMove(newX, m.Pacman.Position.Y) {
				m.Pacman.Position.X = newX
				m.Pacman.Move = utils.Direction{X: 1, Y: 0}
			}
		}

		// Check for dot collection
		for i := len(m.Dots) - 1; i >= 0; i-- {
			if m.Pacman.Position == m.Dots[i].Position {
				m.LevelScore++
				m.Maze[m.Pacman.Position.Y] = utils.ReplaceAtIndex(m.Maze[m.Pacman.Position.Y], ' ', m.Pacman.Position.X) // Replace dot with a space
				m.Dots = append(m.Dots[:i], m.Dots[i+1:]...)
				go m.PlaySound(sound.CHOMP)
				break
			}
		}

		// Check for win condition
		if len(m.Dots) == 0 {
			m.LevelWin = true
			m.GameScore += m.LevelScore
			go m.PlaySound(sound.INTERMISSION)
			return m, nil
		}

		// Check for energizer collection
		for i := len(m.Energizers) - 1; i >= 0; i-- {
			if m.Pacman.Position == m.Energizers[i].Position {
				m.Maze[m.Pacman.Position.Y] = utils.ReplaceAtIndex(m.Maze[m.Pacman.Position.Y], ' ', m.Pacman.Position.X) // Replace dot with a space
				m.Energizers = append(m.Energizers[:i], m.Energizers[i+1:]...)

				// Activate rampant mode
				m.Pacman.RampantState = true
				ghostsEaten = 0
				go m.PlaySound(sound.EATFRUIT)
				return m, m.startRampantTimer()
			}
		}
		cmd := m.checkGhostCollisions()
		return m, cmd
	case ghostMoveMsg:
		// Move ghosts
		for name, g := range m.Ghosts {
			if g.Dead {
				continue
			}
			if m.Pacman.RampantState {
				g.Position = m.escapeMove(g.Position)
			} else {
				switch g.Name {
				case "Blinky":
					g.Position = m.straitMove(g.Position)
				case "Inky":
					g.Position = m.chaosMove(g.Position)
				case "Pinky":
					g.Position = m.predictMove(g.Position)
				case "Clyde":
					g.Position = m.cagyMove(g.Position)
				}
			}
			m.Ghosts[name] = g
		}
		// Start the next tick for ghost movement
		return m, tea.Batch(m.ghostMoveTick(), m.checkGhostCollisions())

	case pacmanBlinkMsg:
		// Toggle the blink state
		m.Pacman.ChewState = !m.Pacman.ChewState
		// Schedule the next blink
		return m, m.pacmanBlinkTick()

	case rampantEndMsg:
		// Start cooldown
		m.Pacman.CooldownState = true
		return m, m.startCooldownTimer()

	case cooldownEndMsg:
		// End cooldown and fully reset Pac-Man's state
		m.Pacman.RampantState = false
		m.Pacman.CooldownState = false
		return m, nil

	case ghostReviveMsg:
		// Revive the ghost at its revival point
		ghost := m.Ghosts[msg.ghostName]
		ghost.Dead = false
		ghost.Position = ghost.RevivalPoint
		m.Ghosts[msg.ghostName] = ghost
		return m, nil
	}

	return m, nil
}

func (m *Model) recordLevelElapsedTime() {
	elapsedTime := int(time.Since(m.CurrentSart).Seconds())
	if m.State.ElapsedTime[m.LevelName] == 0 || elapsedTime < m.State.ElapsedTime[m.LevelName] {
		m.State.ElapsedTime[m.LevelName] = elapsedTime
	}
}
func (m *Model) recordGameScore() {
	if m.State.HighScore < m.LevelScore {
		m.State.HighScore = m.LevelScore
	}
}

func (m *Model) PlaySound(name string) {
	if s, ok := m.Sounds[name]; ok && !m.Mute {
		sound.Play(s)
	}
}
