package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep/speaker"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver {
		if m.lives > 1 {
			// Wait for spacebar to restart the current level
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					lives := m.lives - 1
					m.cancel()                                         // Deduct a life
					newModel := initialModel(m.Config, m.currentLevel) // Restart current level
					newModel.lives = lives                             // Preserve remaining lives
					return newModel, newModel.Init()
				case "q", "ctrl+c":
					return m, tea.Quit
				}
			}
			return m, nil
		}
		// No lives left, offer to restart the game
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				m.cancel()
				newModel := initialModel(m.Config, 0)
				return newModel, newModel.Init()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}
	if m.winGame {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				m.cancel()
				newModel := initialModel(m.Config, 0)
				newModel.winGame = false // Reset the winGame flag
				// Start the timer for ghost movement and blinking
				return newModel, newModel.Init()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		// Stop scheduling commands when the game is won
		return m, nil
	}
	if m.win {
		if m.currentLevel >= len(m.Levels)-1 {
			m.winGame = true
			return m, nil
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				if m.win {
					if m.currentLevel < len(m.Levels)-1 {
						m.currentLevel++
					}
					lives := m.lives
					m.cancel()
					newModel := initialModel(m.Config, m.currentLevel)
					newModel.lives = lives
					// Start the timer for ghost movement and blinking
					return newModel, newModel.Init()
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		speaker.Clear()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.canMove(m.player.position.x, m.player.position.y-1) {
				m.player.position.y--
				m.player.move = direction{x: 0, y: -1}
			}
		case "down":
			if m.canMove(m.player.position.x, m.player.position.y+1) {
				m.player.position.y++
				m.player.move = direction{x: 0, y: 1}
			}
		case "left":
			newX := m.tunnelMove(m.player.position.x - 1)
			if m.canMove(newX, m.player.position.y) {
				m.player.position.x = newX
				m.player.move = direction{x: -1, y: 0}
			}
		case "right":
			newX := m.tunnelMove(m.player.position.x + 1)
			if m.canMove(newX, m.player.position.y) {
				m.player.position.x = newX
				m.player.move = direction{x: 1, y: 0}
			}
		}

		// Check for dot collection
		for i := len(m.dots) - 1; i >= 0; i-- {
			if m.player.position == m.dots[i].position {
				m.score++
				m.maze[m.player.position.y] = replaceAtIndex(m.maze[m.player.position.y], ' ', m.player.position.x) // Replace dot with a space
				m.dots = append(m.dots[:i], m.dots[i+1:]...)
				go m.playSound(SOUND_CHOMP)
				break
			}
		}

		// Check for win condition
		if len(m.dots) == 0 {
			m.win = true
			go m.playSound(SOUND_INTERMISSION)
			return m, nil
		}

		// Check for energizer collection
		for i := len(m.energizers) - 1; i >= 0; i-- {
			if m.player.position == m.energizers[i].position {
				m.score++
				m.maze[m.player.position.y] = replaceAtIndex(m.maze[m.player.position.y], ' ', m.player.position.x) // Replace dot with a space
				m.energizers = append(m.energizers[:i], m.energizers[i+1:]...)

				// Activate rampant mode
				m.player.rampantState = true
				go m.playSound(SOUND_EATFRUIT)
				return m, m.startRampantTimer()
			}
		}
		cmd := m.checkGhostCollisions()
		return m, cmd
	case ghostMoveMsg:
		// Move ghosts
		for name, g := range m.ghosts {
			if g.dead {
				continue
			}
			if m.player.rampantState {
				g.position = m.escapeMove(g.position)
			} else {
				switch g.name {
				case "Blinky":
					g.position = m.straitMove(g.position)
				case "Inky":
					g.position = m.chaosMove(g.position)
				case "Pinky":
					g.position = m.predictMove(g.position)
				case "Clyde":
					g.position = m.cagyMove(g.position)
				}
			}
			m.ghosts[name] = g
		}
		// Start the next tick for ghost movement
		return m, tea.Batch(m.ghostMoveTick(), m.checkGhostCollisions())

	case playerBlinkMsg:
		// Toggle the blink state
		m.player.chewState = !m.player.chewState
		// Schedule the next blink
		return m, m.playerBlinkTick()

	case rampantEndMsg:
		// Start cooldown
		m.player.cooldownState = true
		return m, m.startCooldownTimer()

	case cooldownEndMsg:
		// End cooldown and fully reset Pac-Man's state
		m.player.rampantState = false
		m.player.cooldownState = false
		return m, nil

	case ghostReviveMsg:
		// Revive the ghost at its revival point
		ghost := m.ghosts[msg.ghostName]
		ghost.dead = false
		ghost.position = ghost.revivalPoint
		m.ghosts[msg.ghostName] = ghost
		return m, nil
	}

	return m, nil
}
